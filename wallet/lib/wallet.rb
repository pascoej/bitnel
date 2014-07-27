require '../lib/wallet'
require '../lib/db_models'
require '../lib/wallet_models'
require 'bitcoin'

API_SERVER_URI = "http;//localhost:8080"
API_SERVER_DEPOSIT_PATH = "wallet/notify/deposit"
class Wallet
  attr_reader :currency_data,:node
  def initialize(currency_data)
    Bitcoin.network = currency_data.network
    @currency_data = currency_data
    node_config = NODE_CONFIGS[currency_data.db_currency[:symbol]]
    @node = Bitcoin::Network::Node.new(node_config)
    @handler = WalletHandler.new(self)
  end
  def start
    @node.run
  end
  def newAddress(account = nil)
    key = Bitcoin::generate_key
    privkey = key[0]
    pubkey = key[1]
    address = getAddress(pubkey)
    date = Time.now
    db_address = Address.new(:private_key => privkey, :public_key => pubkey, :address => address, :date => date, :currency => @currency_data.db_currency, :account_uuid => account)
    if (!db_address.save)
      db_address.errors.each { |e| puts e }
      raise 'error generating address!!!!'
    end
    return db_address
  end
  def getAddress pubkey
    hash160 = Bitcoin::hash160(pubkey)
    addr = Bitcoin::encode_address(hash160, @currency_data.addr_version)
    return addr
  end
end

class WalletHandler
  def initialize wallet
    @wallet = wallet
    @node = wallet.node
    @currency_data = wallet.currency_data
    @db_currency = @currency_data.db_currency
    @confirms = @db_currency[:confirms]
    registerListeners
  end
  def registerListeners
    doMissedOutputs
    #Register initial outputs as they are broadcasted...
    @node.subscribe(:tx) do |tx, conf|
      tx.out.each.each_with_index do |out, idx|
        onOutput(tx,out,idx,conf)
      end
    end
    #Register outputs as they come in blocks... (confirms)
    if @confirms > 0
      @node.subscribe(:block) do |block, depth|
        block = @node.store.get_block_by_depth(depth - conf + 1)
        next  unless block
        block.tx.each do |tx|
          tx.out.each.with_index do |out, idx|
            onOutput(tx,out,idx,conf)
          end
        end
      end
    end
  end
  def doMissedOutputs
    return unless @db_currency[:last_tx] != nil
    return unless @db_currency[:last_idx] != nil
    last_hash = @db_currency[:last_tx]
    last_idx = @db_currency[:last_idx].to_i

    return unless last_tx = @node.store.get_tx(last_hash)
    return unless last_idx = last_tx.out[last_idx]

    depth = @node.store.get_depth
    (last_tx.get_block.depth..depth).each do |i|
      blk = @node.store.get_block_by_depth(i)
      blk.tx.each do |tx|
        tx.out.each.with_index do |out, idx|
            conf = (depth - blk.depth + 1)
            onOutput(tx, out, idx,conf)
        end
      end
    end
  end
  def onOutput tx, out, idx, conf
    puts "swag swag like calliou"
    script = Bitcoin::Script.new(out.pk_script)
    address = script.get_address
    value = out.value
    @db_currency = Currency.first(:id => @db_currency[:id])
    record = Address.first(:address => address)
    if record != nil
      puts "swag"
      onTransaction(address,value,tx,out,idx,conf)
    end
  end
  def onTransaction addr, value, tx, out, idx, conf
    confirmed = conf > @confirms
    puts "too much swag"
    @db_currency = Currency.first(:id => @db_currency[:id])
    db_out = Output.first_or_create(:tx_hash => tx.hash, :address => addr, :value => value, :idx => idx, :currency => @db_currency, :date => Time.now)
    if (!db_out.save)
      db_out.errors.each { |e| puts e.inspect }
    end
    db_out = Output.get(id)
    db_out.update(:value => value, :tx_hash => tx.hash, :confirmed => confirmed, idx => idx)
    bodyHash = {"address" => addr[:address], "confirmed" => confirmed, "tx_hash" => tx.hash, "value" => value}
    body = bodyHash.to_json.to_s
    request = Request.create(:type => "deposit", :done => false, :uri => API_SERVER_URI, :path => API_SERVER_DEPOSIT_PATH, :body => body)
    if (!request.save)
      request.errors.each { |e| puts e.inspect }
    end
  end
end




NODE_CONFIGS = {
    "btc" =>  {
        :network => :bitcoin,
        :listen => ["0.0.0.0", 7000],
        :connect => [],
        :command => ["127.0.0.1", 6000],
        :storage => "utxo::sqlite://~/.bitcoin-ruby/<network>/blocks.db",
        :announce => false,
        :external_port => nil,
        :mode => :full,
        :cache_head => true,
        :index_nhash => false,
        :index_p2sh_type => false,
        :dns => true,
        :epoll_limit => 10000,
        :epoll_user => nil,
        :addr_file => "~/.bitcoin-ruby/<network>/peers.json",
        :log => {
            :network => :info,
            :storage => :info,
        },
        :max => {
            :connections_out => 8,
            :connections_in => 32,
            :connections => 8,
            :addr => 256,
            :queue => 501,
            :inv => 501,
            :inv_cache => 0,
            :unconfirmed => 100,
        },
        :intervals => {
            :queue => 1,
            :inv_queue => 1,
            :addrs => 5,
            :connect => 5,
            :relay => 0,
        },
        :import => nil,
        :skip_validation => false,
        :check_blocks => 1000,
    },
    "rgt" => {
    :network => :regtest,
    :listen => ["0.0.0.0", 7001],
    :connect => ["127.0.0.1",8333],
    :command => ["127.0.0.1", 6001],
    :storage => "utxo::sqlite://~/.bitcoin-ruby/<network>/blocks.db",
    :announce => false,
    :external_port => nil,
    :mode => :full,
    :cache_head => true,
    :index_nhash => false,
    :index_p2sh_type => false,
    :dns => true,
    :epoll_limit => 10000,
    :epoll_user => nil,
    :addr_file => "~/.bitcoin-ruby/<network>/peers.json",
    :log => {
        :network => :info,
        :storage => :info,
    },
    :max => {
        :connections_out => 8,
        :connections_in => 32,
        :connections => 8,
        :addr => 256,
        :queue => 501,
        :inv => 501,
        :inv_cache => 0,
        :unconfirmed => 100,
    },
    :intervals => {
        :queue => 1,
        :inv_queue => 1,
        :addrs => 5,
        :connect => 5,
        :relay => 0,
    },
    :import => nil,
    :skip_validation => false,
    :check_blocks => 1000,
}
}
require 'sinatra'

get '/' do
  'Welcome to Bitnel Wallet API!'
end

post '/new_address' do
  account_uuid = params[:account_uuid]
  currency_symbol = params[:currency_symbol].to_s
  wallet = $wallet_server.wallets[currency_symbol]
  if (wallet == nil)
    status 500
    return
  end
  addr = wallet.newAddress(account_uuid)
  addr[:address]
end

require '../lib/wallet'
require '../lib/db_models'
require '../lib/wallet_models'
require '../lib/api'

class WalletServer
  attr_accessor :wallets
  def initialize
    db_currency = Currency.first(:symbol=>"rgt")
    currency_data = CurrencyData.new(db_currency)
    #db_btc_currency = Currency.first(:symbol=>"btc")
    #btc_data = CurrencyData.new(db_btc_currency)
    #kick ali out of call pls
    @wallets = {"rgt" => Wallet.new(currency_data)}sou
  end
  def start
    @wallets.each_value do |wallet|
      Thread.new do
        wallet.start
      end
    end
  end
end

$wallet_server = WalletServer.new
$wallet_server.start
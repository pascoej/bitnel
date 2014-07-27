require 'data_mapper'
#DataMapper::Logger.new($stdout, :debug)

DataMapper.setup(:default, 'mysql://root:test@localhost/bn_wallet')


class Currency
  include DataMapper::Resource

  property :id,           Serial
  property :currency_id,  Integer, :required => true, :key => true
  property :network,      String
  property :symbol,       String, :required => true, :key => true
  property :last_tx,      String,  :length => 64
  property :last_idx,     Integer
  property :confirms,     Integer
end
class Address
  include DataMapper::Resource

  property :id,           Serial
  property :account_uuid,      String, :length => 60
  property :address,      String, :length => 80, :key => true
  property :private_key,  String, :length => 80
  property :public_key,   String, :length => 140
  property :date,         DateTime
  belongs_to :currency
end
class Output
  include DataMapper::Resource

  property :id,           Serial
  property :tx_hash,      String, :length => 64
  property :confirmed,    Boolean, :default => false
  property :value,        Integer, :min => 0, :max => 18_446_744_073_709_551_615
  property :idx,          Integer
  property :date,         DateTime
  property :spent,         Boolean, :default => false

  belongs_to :address
  belongs_to :currency
end


class Request
  include DataMapper::Resource

  property :id,           Serial
  property :type,         String
  property :done,         Boolean, :default => false
  property :uri,          String
  property :path,         String
  property :body,         Json
  property :date,         DateTime
end

class ColdWalletAddress
  include DataMapper::Resource

  property :id,           Serial
  property :address,      String, :length => 80

  belongs_to :currency
end

class Withdraw
  include DataMapper::Resource

  property :id,           Serial
  property :account,      String, :length=>64
  property :nonce,        Integer, :key => true
  property :amount,       Integer, :min => 0, :max => 18_446_744_073_709_551_615
  property :fee,          Integer
  property :address,      String, :length => 80
  property :done,         Boolean, :default => false
  property :date,         DateTime
  has 1,    :WithdrawTransaction
  belongs_to :currency
end

class WithdrawTransaction
  include DataMapper::Resource

  property :id,           Serial
  property :json_tx,      Json
  property :accepted,     Boolean, :default => false
  property :time,         DateTime
  belongs_to :withdraw
end
DataMapper.finalize
DataMapper.auto_migrate!
rgt = Currency.first_or_create(:currency_id => 0, :confirms => 3, :symbol => "rgt", :network => "regtest")
btc = Currency.first_or_create(:currency_id => 0, :confirms => 3, :symbol => "btc", :network => "bitcoin")

puts rgt.save
puts btc.save
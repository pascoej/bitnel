require 'rubygems'
require 'bitcoin'

class CurrencyData
  attr_reader :addr_version, :network, :network_info, :db_currency
  def initialize db_currency
    @db_currency = db_currency
    @network = db_currency[:network].to_sym
    @network_info = Bitcoin::NETWORKS[@network]
    @addr_version = @network_info[:address_version]
  end
end

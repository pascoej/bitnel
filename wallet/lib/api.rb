require 'rubygems'
require 'sinatra'
get '/' do
  'Welcome to Bitnel API!'
  p @wallet_server.inspect
end
post '/new_address' do
  account_uuid = params[:account_uuid]
  currency_symbol = params[:currency_symbol].to_s
  wallet = @wallet_server.wallets[currency_symbol]
  if (wallet == nil)
    status HTTPForbidden
  end
  addr = wallet.newAddress(account_uuid)
  addr[:address]
end
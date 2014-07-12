require 'sinatra'
require 'bitcoin'
require 'active_support'
j = ActiveSupport::JSON
get '/new-address' do
	key = Bitcoin::generate_key
	return j.encode(key)
end;
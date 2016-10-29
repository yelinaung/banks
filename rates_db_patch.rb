require 'rethinkdb'
include RethinkDB::Shortcuts
r.connect(host: 'localhost', port: 28_015).repl
cursor = r.db('test').table('currency').run
count = 0
cursor.each do |d|
  p "Updating document #{d['id']}"
  orig_rates = d['rates']
  next if orig_rates[0]['USD'].nil? # HACK: for mixed tables
  new_rates = []
  orig_rates.each do |rate|
    rate.each do |k, v|
      h = v.clone
      h['currency_shortcode'] = k
      new_rates.push h
    end
  end
  p new_rates
  r.db('test').table('currency').update(rates: new_rates).run
  count += 1
  p '======='
end
puts "Total #{count} documents updated"

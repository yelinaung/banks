require 'rethinkdb'
include RethinkDB::Shortcuts
r.connect(host: 'localhost', port: 28_015).repl
cursor = r.db('test').table('currency').run
cursor.each do |d|
  p d
  orig_rates = d['rates']
  next if orig_rates[0]['USD'].nil? # HACK: for mixed tables
  new_rates = []
  orig_rates.each do |rate|
    p rate
    rate.each do |k, v|
      h = v.clone
      h['currency_shortcode'] = k
      p h
      p v
      new_rates.push h
    end
  end
  p new_rates
  r.db('test_cloned').table('currency_sample').update(rates: new_rates).run
  p '======='
end

#! /usr/bin/env ruby
require 'optparse'
require 'tempfile'

options = {
  before: 'master',
  after: 'HEAD',
  bench_time: '1s'
}

OptionParser.new do |opts|
  opts.banner = 'Runs benchmarks on two branches and compares the results'
  opts.on('-b hash', '--before hash') do |v|
    options[:before] = v
  end

  opts.on('-a hash', '--after hash') do |v|
    options[:after] = v
  end

  opts.on('-t time', '--bench_time time') do |v|
    options[:bench_time] = v
  end
end.parse!

benchmark_options = "-run '^$' -bench '.' -benchmem -benchtime #{options[:bench_time]}"
return_to = `git rev-parse --abbrev-ref HEAD`
before_hash = `git rev-parse #{options[:before]}`.strip
after_hash = `git rev-parse #{options[:after]}`.strip

bf = Tempfile.new('before')
af = Tempfile.new('after')
begin
  `git checkout #{before_hash} 2>&1`

  puts "benchmarking #{before_hash}"
  bf.write `go test #{benchmark_options} ./...`

  `git checkout #{after_hash} 2>&1`
  puts "benchmarking #{after_hash}"
  af.write `go test #{benchmark_options} ./...`
  af.close
  bf.close

  `go get golang.org/x/tools/cmd/benchcmp`
  comparison = `$GOPATH/bin/benchcmp #{bf.path} #{af.path}`

  puts RUBY_PLATFORM
  puts comparison
rescue StandardError => e
  puts e
ensure
  bf.unlink
  af.unlink
  `git checkout #{return_to} 2>&1`
end

#!/usr/bin/env ruby

require 'rainbow'

expected = File.new("assets/nestest.log").readlines.map { |l| l.gsub("\r\n", '')[0..72] }
actual = `go run src/nestest.go`.split /\n/

actual.zip(expected).each_with_index do |(a, e), i|
  if e.nil?
    break
  elsif e != a
    expected_diff = e.chars.zip(a.chars).map { |ec, ac|
      if ac != ec
        ec.foreground(:green)
      else
        ec
      end
    }.join ''

    actual_diff = a.chars.zip(e.chars).map { |ac, ec|
      if ac != ec
        ac.foreground(:red)
      else
        ac
      end
    }.join ''

    if ![
      5013, 5143, 5175, 5176, 5177, 5178, 5179, 5180, # Problems with $A9A9 -- is this a device?
      8981, 8983, 8985, 8987, 8989 # Problems with $40xx, this is the APU, so no wonder there
    ].include?(i+1)
      ([0, i-15].max..i-1).each { |i| puts "                " + actual[i] }

      puts sprintf("%-5d expected: %s", i + 1, expected_diff)
      puts sprintf("%-5d actual  : %s", i + 1, actual_diff)
      puts "\r\n"

      break
    end
  end
end


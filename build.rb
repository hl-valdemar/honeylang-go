#!/usr/bin/ruby

is_release = false

for arg in ARGV do
  case arg
  when "--release"
    is_release = true
  else
    abort "Unrecognized argument: #{arg}"
  end
end

if is_release then
  success = system "go build -ldflags '-s -w'"
  if !success then
    puts "Failed to build honeylang (release)"
  end
else
  success = system "go build ."
  if !success then
    puts "Failed to build honeylang (debug)"
  end
end


#!/usr/bin/ruby

# parse arguments
is_release = false

for arg in ARGV do
  case arg
  when "--release"
    is_release = true
  else
    abort "Unrecognized argument: #{arg}"
  end
end

# generate enum string methods
system "go generate ./lexer"

# build the honey compiler
if is_release then
  success = system "go build -ldflags '-s -w'"
  if !success then
    abort "Failed to build honeylang (release)"
  end
else
  success = system "go build ."
  if !success then
    abort "Failed to build honeylang (debug)"
  end
end


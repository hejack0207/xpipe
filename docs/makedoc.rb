#!/usr/bin/env ruby
#
#   Scrape document strings from the GO files and convert them into
#   a single Markdown document.
#

def process_file(filename)
    docs = []
    currDoc = nil

    File.open(filename, "r") do |f|
        f.each_line do |line|
            if (line =~ /^\/\/\/\//) then
                syntax = $'.strip()
                name, _ = syntax.split(" ", 2)

                currDoc = { :name => name, :syntax => syntax, :lines => "" }
                docs << currDoc
            elsif ((line =~ /^\/\//) and (currDoc != nil)) then
                currDoc[:lines] += $'.strip() + "\n"
            else
                currDoc = nil
            end
        end
    end

    docs
end

def write_md_file(docs)
    puts "# Processes"
    puts

    docs.each do |d|
        puts "### #{d[:name]}"
        puts
        puts "```"
        puts d[:syntax]
        puts "```"
        puts
        puts d[:lines].strip()
        puts
    end
end

allDocs = []

Dir.glob("src/**/*.go") do |file|
    allDocs += process_file(file)
end

allDocs.sort! { |x, y| x[:name] <=> y[:name] }

write_md_file(allDocs)

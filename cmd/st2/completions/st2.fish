complete st2 -f
complete st2 -r -f -s r -l root -d 'The root struct name (default: Root)'
complete st2 -r -F -s i -l input -d 'Input file, if not set, it will read from stdio'
complete st2 -l rc -d 'Read input from clipboard'
complete st2 -r -f -s s -l src -a "json yaml proto thrift go csv" -d 'The source data type, it will use the suffix of the input file if not set'
complete st2 -r -f -s d -l dst -a "go proto thrift" -d 'The destination data type, it will use the suffix of the output file if not set'
complete st2 -r -F -s o -l output -d 'Output file, if not set, it will write to stdout'
complete st2 -l wc -d 'Write output to clipboard'

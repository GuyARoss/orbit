echo "\n"

for file in examples/*; do
    rm -rf "$file"/orbit

    cp orbit "$file"/orbit

    echo "linked $file"
done

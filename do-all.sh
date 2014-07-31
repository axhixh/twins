rm data/*.txt data/*.cln data/*.tf
for v in data/*.xml; do ./to-txt $v; done
for v in data/*.txt; do ./clean $v; done
for v in data/*.cln; do ./tf $v;done
./global-idf data
for v in data/*.tf; do ./idf -input $v;done
./similar data

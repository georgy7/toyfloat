import matplotlib.pyplot as plt
import csv

values = []
precision = []

tsv_file = open("precision.tsv")
read_tsv = csv.reader(tsv_file, delimiter="\t")

for row in read_tsv:
    v = float(row[0])
    p = float(row[1])

    if v < 20000:
        values.append(v)
        precision.append(p)

figure = plt.figure()

figure.patch.set_facecolor('#e6e6e6')
plt.axes().set_facecolor('#eeeeee')

figure.set_size_inches(5, 5)
plt.subplots_adjust(left=0.162, bottom=0.129, right=0.954, top=0.954)

plt.plot(values, precision, '.')

plt.xlabel('Floating Point Values')
plt.ylabel('Floating Point Precision')

plt.xscale('log')
plt.yscale('log')

plt.grid(True, color='#c7c7c7')

# plt.title('Precision of Toyfloat')

# plt.show()
plt.savefig('precision.png')

import matplotlib.pyplot as plt
import csv

values = []
precision = []

tsv_file = open("precision.tsv")
read_tsv = csv.reader(tsv_file, delimiter="\t")

for row in read_tsv:
    v = float(row[0])
    p = float(row[1])

    if (v < 20000) and (p > 0):
        values.append(v)
        precision.append(p)

figure = plt.figure()

bg_color = '#e6e6e6'
figure.patch.set_facecolor(bg_color)
plt.axes().set_facecolor(bg_color)

figure.set_size_inches(5, 5)
plt.subplots_adjust(left=0.163, bottom=0.138, right=0.931, top=0.907)

plt.plot(values, precision, '.')

plt.xlabel('Floating Point Values')
plt.ylabel('Floating Point Precision')

plt.xscale('log')
plt.yscale('log')

plt.title('Precision of Toyfloat')

# plt.show()
plt.savefig('precision.png')

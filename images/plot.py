import matplotlib.pyplot as plt

values = [1,2,3]
precision = [2,4,1]

plt.plot(values, precision)

plt.xlabel('Floating Point Values')
plt.ylabel('Floating Point Precision')

plt.title('Precision of Toyfloat')

plt.show()

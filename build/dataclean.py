# 说明
# 该文件可实现数据清洗。
# 输入：源list
# 输出：清洗的list
# 数据结构为list

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import sys

#定义查找函数，找两边
def get_leftandright(lst, index):
    if 0 < index < len(lst) - 1:
        return (lst[index - 1]+lst[index + 1])/2
    else:
        return None

#找左边
def get_left(lst, index):
    if index>=3:
        return np.mean(lst[index-3:index])
    else:
        return None

def data_clean(data):
    positions = [i for i, x in enumerate(data) if x is None]
    print("input positions:",positions)
    if len(positions)>0 and positions[0]>0 and positions[len(positions)-1]<len(data)-1:
        for x in positions:
            data[x]=get_leftandright(data,x)

    lower=np.mean(data)-2*np.std(data);
    upper=np.mean(data)+2*np.std(data);
    x=[num for num in data if num > upper or num<lower]
    positions = [index for index, value in enumerate(data) if value > upper or value<lower]
    if positions[0]>2:
        for x in positions:
            data[x]=get_left(data,x)

    return data

#命令行参数中读取输入文件和输出文件名
datafile=sys.argv[1]
resultfile=sys.argv[2]

#读取输入数据
df = pd.read_csv(datafile, delimiter=",", header=None)  # 如果你的数据是以tab分隔的
data = df.iloc[:, 7].values.tolist()
print("input data:",data)

#数据清洗
data=data_clean(data)

#数据写回原始读入的list中
for i in range(len(data)):
    df.iloc[i,7]=data[i]

#将清洗后的df数据写入文件
df.to_csv(resultfile, header=False, index=False)

print('0')







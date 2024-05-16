# 说明：
# 该文件是预测程序
import sys
import numpy as np
from keras.models import Sequential
from keras.layers import Dense, Activation, Dropout, LSTM
import math
from sklearn.metrics import mean_squared_error
from sklearn.preprocessing import MinMaxScaler
import logging
import pandas as pd
from datetime import datetime, timedelta

logging.disable(logging.DEBUG)
logging.disable(logging.WARNING)
# 定义create_dataset()函数，构建数据集合
def create_dataset(dataset, look_back=1):
    dataX, dataY = [], []
    for i in range(len(dataset) - look_back-1):
        a = dataset[i: (i+look_back)]
        dataX.append(a)
        dataY.append(dataset[i+look_back])
    return np.array(dataX), np.array(dataY)

# 定义tem_predict()函数，建立预测程序.
def train_model(dataset):
    #归一化操作，将数据标准化到0到1
    scaler = MinMaxScaler(feature_range=(0, 1))
    dataset = scaler.fit_transform(dataset.reshape(-1, 1))

    train_size = int(len(dataset)*0.95)
    # test_size = len(dataset)-train_size
    train, test = dataset[0: train_size], dataset[train_size: len(dataset)]
    #print('Shape of array train:', train.shape)

    look_back = 1 #步长1 
    trainX, trainY = create_dataset(train, look_back)  #训练集
    testX, testY = create_dataset(test, look_back)     #测试集

    model = Sequential()#顺序模型
    model.add(LSTM(input_shape=(None,1),units=100, return_sequences=False))
    model.add(Dense(units=1))
    model.add(Activation('linear'))
    model.compile(loss='mean_squared_error', optimizer='Adam')
    model.summary()

    #verbose=0，不打印输出。
    #validation_split=0.1，样本10%作为验证集
    history = model.fit(trainX, trainY, batch_size=64, epochs=50,validation_split=0.1, verbose=0)

    return model

def predict(model, dataset):
    #归一化操作，将数据标准化到0到1
    scaler = MinMaxScaler(feature_range=(0, 1))
    dataset = scaler.fit_transform(dataset.reshape(-1, 1))
    #print(dataset)
    dataset=np.array(dataset)
    #print('Shape of array b:', dataset.shape)
    #print("dataset:",dataset)
    # 执行预测。
    predictResult = model.predict(dataset)
 
    # 反归一化，获取真实值。
    predictResult = scaler.inverse_transform(predictResult)
    return predictResult

def getNexDateStr(dateStr):
    date = datetime.strptime(dateStr, "%Y-%m-%d 00:00:00")
    next_date = date + timedelta(days=1)
    return next_date.strftime("%Y-%m-%d 00:00:00")

#命令行参数中读取输入文件和输出文件名
hisfile=sys.argv[1]
resultfile=sys.argv[2]

#读取历史数据
df = pd.read_csv(hisfile, delimiter=",", header=None)  # 如果你的数据是以tab分隔的
data = df.iloc[:, 5].values.tolist()
data=np.diff(data,0)
dataset=np.array(data)
#print("dataset:",dataset)
model=train_model(dataset)
last_day=dataset[-1]
#预测未来7天的温度值
#定义一个列表，存放预测结果
predictResultList=[]
for fd in range(7):
    print(last_day)
    predictResult=predict(model,last_day)
    print(predictResult)
    last_day=predictResult
    predictResultList.append(predictResult[0][0])   

# 读取最后一行
last_row = df.iloc[-1]
#print(last_row)
# 取出其中的第一列的值
last_day = last_row[3]
#print(last_day)
device_type_id=last_row[0]
device_id=last_row[1]
sensor_id=last_row[2]
time=last_row[4]

#将预测结果写入文件
with open(resultfile, 'w') as f:
    for i in range(len(predictResultList)):
        #取last_day的下一天
        last_day=getNexDateStr(last_day)
        f.write(str(device_type_id) + ',' + str(device_id) + ',' + str(sensor_id) + ','+ str(last_day)+',' + str(time) + ',' + str(predictResultList[i]) + '\n')
    f.close()

print('0')

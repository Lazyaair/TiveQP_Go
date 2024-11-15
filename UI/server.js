const express = require('express');
const fs = require('fs');
const path = require('path');
const app = express();
const port = 3000;

// 允许跨域请求
app.use((req, res, next) => {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    next();
});

// 根路径路由
app.get('/', (req, res) => {
    res.send('欢迎来到 Node.js 服务器！');
});

// 读取文本文件，解析并返回数据
app.get('/locations', (req, res) => {
    const filePath = 'C:\\Users\\wzy\\Desktop\\TiveQP-main\\TiveQP\\ServeDB\\Data\\20k_random.txt'; 
    fs.readFile(filePath, 'utf8', (err, data) => {
        if (err) {
            return res.status(500).send('读取文件失败');
        }
        
        const locations = [];
        const lines = data.split('\n');
        
        lines.forEach(line => {
            const parts = line.split('**');
            if (parts.length === 8) {
                const [type, city, lat, lng, openHour, openMin, closeHour, closeMin] = parts;
                locations.push({
                    type,
                    city,
                    lat: parseFloat(lat),
                    lng: parseFloat(lng),
                    openHour: parseInt(openHour),
                    openMin: parseInt(openMin),
                    closeHour: parseInt(closeHour),
                    closeMin: parseInt(closeMin)
                });
            }
        });
        
        res.json(locations); // 返回解析后的数据
    });
});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});

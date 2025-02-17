**This project outputs the mini 24-hour ticker information for the BTC/USDT pair from Binance. It includes detailed data such as**
**the last 24-hour price change, price information, volume, and the number of trades.**

##  How It Works
**_docker-compose up --build_ command will be enough to run the project.**
**You can view the metrics as text on ```localhost:8080``` and access the metrics visually on ```localhost:9090```**

```bash
    w: Weighted Average Price
    B: Best Bid Quantity
    e: Event Type
    E: Event Time
    P: Price Change Percent
    b: Best Bid Price
    A: Best Ask Quantity
    o: Open Price
    F: First Trade ID
    L: Last Trade ID
    x: Previous Close Price
    a: Best Ask Price
    h: High Price
    C: Statistics Close Time
    n: Total Number of Trades
    s: Symbol
    p: Price Change (Absolute)
    c: Current Price (Close Price)
    Q: Last Trade Quantity
    l: Low Price
    v: Volume (Base Asset Volume)
    q: Quote Asset Volume
    O: Statistics Open Time
```
一个简单的将长链接转换为短链接的后台系统。

压力测试结果如下：

使用siege工具对后台进行测试，使用的命令为:  siege -c 200 -r 50 -f searchUrls.txt。

模拟200个用户重复进行50次将长链接转换为短链接的请求，请求内容为searchUrls.txt中的随机url。

测试的结果如下：

当mysql和redis中都不存在数据记录，系统需要对长链接进行转换时：

![情况1压测结果](https://ws1.sinaimg.cn/large/006tNc79ly1g3jgmt2mu0j30ca06tjs3.jpg)

当mysql中存在数据记录，但是redis中缓存过期时：

![情况2压测结果](http://ww4.sinaimg.cn/large/006tNc79ly1g3jgppr4njj30bh06qdgi.jpg)

当直接命中redis缓存时：

![情况3压测结果](http://ww4.sinaimg.cn/large/006tNc79ly1g3jgq8wmxvj30ch06qjs2.jpg)


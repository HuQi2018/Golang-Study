<html>
<head>
    <title></title>
</head>
<body>
{{.}}
<form action="/login" method="post" enctype="multipart/form-data">
    用户名:<input type="text" name="username"><br>
    密码:<input type="password" name="password"><br>
    身份证号:<input type="text" name="usercard"><br>
    邮箱地址:<input type="text" name="email"><br>
    手机号码:<input type="text" name="mobile"><br>
    年龄:<input type="text" name="age"><br>
    真实姓名:<input type="text" name="realname"><br>
    水果：
    <select name="fruit">
        <option value="apple">apple</option>
        <option value="pear">pear</option>
        <option value="banane">banane</option>
    </select><br>
    复选框：
    <input type="checkbox" name="interest" value="football">足球
    <input type="checkbox" name="interest" value="basketball">篮球
    <input type="checkbox" name="interest" value="tennis">网球<br>
    <input type="submit" value="提交"><br>
</form>
</body>
</html>

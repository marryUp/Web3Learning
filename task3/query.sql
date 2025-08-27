/**1.基本sqlCRUD操作
假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
*/

-- 1.1编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
-- insert into students (name, age, grade) values ('张三', 20, '三年级');

-- 1.2编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
-- select * from students where age > 0;

-- 1.3编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
-- update students set grade = '四年级' where name = '张三';

-- 1.4编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
-- delete FROM students where age < 15;


/**
2.编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

-- select * from accounts;
-- insert into accounts (Id, balance) values (1, 1000);
-- insert into accounts (Id, balance) values (2, 500);

--创建事务
BEGIN TRANSACTION;
-- 检查账户 A 的余额是否足够
update accounts set balance = balance - 100 where Id = 1 and balance>100;
-- 如果余额足够，并且上面语句已经执行成功，则向账户 B 增加 100 元
update accounts set balance = balance + 100 where Id = 2 AND (SELECT changes() = 1);
--提交事务
COMMIT;

select * from accounts;
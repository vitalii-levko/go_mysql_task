SELECT
 `departments`.`dept_name` AS 'department',
 COUNT(`employees`.`emp_no`) AS 'employee count',
 SUM(`salaries`.`salary`) AS 'sum salary'
FROM `departments`
 INNER JOIN `dept_emp`
  ON `departments`.`dept_no` = `dept_emp`.`dept_no`
 INNER JOIN `employees`
  ON `dept_emp`.`emp_no` = `employees`.`emp_no`
 INNER JOIN `salaries`
  ON `employees`.`emp_no` = `salaries`.`emp_no`
WHERE `dept_emp`.`to_date` > NOW()
 AND `salaries`.`to_date` > NOW()
GROUP BY `departments`.`dept_name`
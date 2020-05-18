SELECT
 `departments`.`dept_name` AS 'department',
 `employees`.`first_name` AS 'first name',
 `employees`.`last_name` AS 'last name',
 `titles`.`title`,
 `salaries`.`salary`
FROM `employees`
 INNER JOIN `dept_manager`
  ON `employees`.`emp_no` = `dept_manager`.`emp_no`
 INNER JOIN `departments`
  ON `dept_manager`.`dept_no` = `departments`.`dept_no`
 INNER JOIN `titles`
  ON `employees`.`emp_no` = `titles`.`emp_no`
 INNER JOIN `salaries`
  ON `employees`.`emp_no` = `salaries`.`emp_no`
WHERE `dept_manager`.`to_date` > NOW()
 AND `titles`.`to_date` > NOW()
 AND `salaries`.`to_date` > NOW()
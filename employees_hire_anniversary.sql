SELECT
 `departments`.`dept_name` AS 'department',
 `titles`.`title`,
 `employees`.`first_name` AS 'first name',
 `employees`.`last_name` AS 'last name',
 `employees`.`hire_date` AS 'hire date',
 TIMESTAMPDIFF(YEAR, `employees`.`hire_date`, NOW()) AS 'work years'
FROM `employees`
 INNER JOIN `dept_emp`
  ON `employees`.`emp_no` = `dept_emp`.`emp_no`
 INNER JOIN `departments`
  ON `dept_emp`.`dept_no` = `departments`.`dept_no`
 INNER JOIN `titles`
  ON `employees`.`emp_no` = `titles`.`emp_no`
WHERE `dept_emp`.`to_date` > NOW()
 AND `titles`.`to_date` > NOW()
 AND MONTH(NOW()) = MONTH(`employees`.`hire_date`)

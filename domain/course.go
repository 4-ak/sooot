package main

type Course struct {
	Name      string
	Professor string
	RList     ReviewList
}

type CourseList struct {
	List []Course
}

func (c *CourseList) Input(prof, cn string) {
	c.List = append(c.List, Course{prof, cn, ReviewList{}})
}

func (c *CourseList) Update(prof, cn string, num int) {
	c.List[num-1].Professor = prof
	c.List[num-1].Name = cn
}

func (c *CourseList) Delete(num int) {
	bb := c.List
	c.List = bb[:num-1]
	bb = bb[num:]
	c.List = append(c.List, bb...)
}

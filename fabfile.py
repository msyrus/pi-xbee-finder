from invoke import task

@task
def build(c):
	c.local("./build.sh", env=c.env)

@task
def deploy(c):
	build(c)
	c.run("mkdir -p app")
	c.run("pkill finder", warn=True)
	c.put("build/finder", "./app/finder")
	c.run("./app/finder")

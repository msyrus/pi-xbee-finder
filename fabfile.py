from invoke import task

@task
def build(c):
	c.local("./build.sh", env=c.env)

@task
def deploy(c):
	c.run("mkdir -p app")
	stop(c)
	c.put("build/finder", "./app/finder")
	c.sudo("./app/finder")

@task
def stop(c):
	c.sudo("pkill -15 finder", warn=True)

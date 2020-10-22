def main(ctx):
    steps = [step(x, ctx.pull_policy, ctx.commands) for x in ctx.tags]

    pipeline = {
        'version': '1',
        'steps': steps,
    }   

    return pipeline
    
def step(tag, pull_policy, commands):
    return {
        "name": "build %s" % tag,
        "image": "golang:%s" % tag,
        "pull": pull_policy,
        'commands': commands.values(),
    }
load("testdata/step/with_module/foo_module.star", "foo")

def main(ctx):
    return {
        'version': '1',
        'steps': [
            {
                "name": "build_%s" % foo.foo,
                "image": "alpine:latest",
                'commands': [
                    "echo %s" % foo.foo
                ]
            }, {
                "name": "build_%s" % foo.bar(),
                "image": "alpine:latest",
                'commands': [
                    "echo %s" % foo.bar()
                ]
            }
        ],
    }

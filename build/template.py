from jinja2 import Environment, FileSystemLoader
import re
import os

templates_env = Environment(
    loader=FileSystemLoader("./"),
    lstrip_blocks=True,
    trim_blocks=True
)

config_template = templates_env.get_template("docker-compose.yml.j2")

with open("./docker-compose.yml", "w") as f:
    refName = os.getenv("GITHUB_REF")
    if not refName:
        raise ValueError("GITHUB_REF not set")
    refName = re.sub(r"/", "_", refName)

    repository = os.getenv("GITHUB_REPOSITORY")
    if not repository:
        raise ValueError("GITHUB_REPOSITORY not set")

    sha = os.getenv("GITHUB_SHA")
    if not sha:
        raise ValueError("GITHUB_SHA not set")

    domain = os.getenv("DOMAIN")
    if not domain:
        raise ValueError("DOMAIN not set")

    f.write(config_template.render(
        {
            "REF_NAME": refName,
            "REPOSITORY": repository,
            "SHA": sha,
            "DOMAIN": domain,
        }))
    f.close()


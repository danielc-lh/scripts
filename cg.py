import json
from pathlib import Path
import subprocess

OUT_DIR = Path(__file__).parent.parent / Path("load-balancer/modules/customer")

# does user have hcl2json installed?
try:
    subprocess.run("hcl2json -h", shell=True, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
except subprocess.CalledProcessError:
    print("hcl2json is not installed. Please see installation instructions at https://github.com/tmccombs/hcl2json?tab=readme-ov-file#installation")

# generate json from tf file containing routing rules
try:
    subprocess.run(f"hcl2json {OUT_DIR}/target-proxies.tf > {OUT_DIR}/target-proxies.tf.json", shell=True, check=True)
except subprocess.CalledProcessError as e:
    print(f"Error occurred: {e}")

with open(OUT_DIR / "target-proxies.tf.json", "r") as f:
    target_proxies_json = f.read()
    t = json.loads(target_proxies_json)

    rules = t["resource"]["google_compute_url_map"]["default"][0]["path_matcher"][0]

    dynamic_path_rules = rules["dynamic"]["path_rule"]
    print(dynamic_path_rules)

    other_rules = rules["path_rule"]
    print(other_rules)

    all_rules = dynamic_path_rules + other_rules

    print(all_rules)

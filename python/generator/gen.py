import json
import uuid
import random
from faker import Faker
from datetime import datetime, timedelta

# Initialize Faker to generate realistic-looking data
fake = Faker()

# --- Helper functions for generating complex patient data ---

def _generate_codes(system, code, display, source="emr"):
    """Generates a list containing a code dictionary."""
    return [{"code": code, "system": f"http://{system}.org", "display": display, "source": source}]

def _generate_medications(num=1):
    """Generates a list of random medication records."""
    medications = []
    med_list = ["Lisinopril 20 MG", "Metformin 500 MG", "Atorvastatin 40 MG", "Amlodipine 10 MG", "Hydrochlorothiazide 25 MG"]
    for _ in range(num):
        med_name = random.choice(med_list)
        medications.append({
            "name": med_name,
            "prescribedDate": fake.date_time_between(start_date='-10y', end_date='-1y').isoformat() + "Z",
            "status": random.choice(["active", "stopped"]),
            "codes": _generate_codes("nlm.nih.gov/research/umls/rxnorm", str(fake.random_number(digits=6)), med_name),
            "encounterId": str(uuid.uuid4()),
            "ingredients": []
        })
    return medications

def _generate_vitals(num=5):
    """Generates a list of random vital signs."""
    vitals = []
    vital_types = {
        "Body Weight": ("29463-7", "kg", (50, 120)),
        "Body Height": ("8302-2", "cm", (150, 200)),
        "Body Mass Index": ("39156-5", "kg/m2", (18, 40))
    }
    for _ in range(num):
        vital_name, (code, unit, val_range) = random.choice(list(vital_types.items()))
        vitals.append({
            "codes": _generate_codes("loinc", code, vital_name),
            "value": str(random.uniform(*val_range)),
            "unit": unit,
            "reportedTime": fake.date_time_between(start_date='-5y', end_date='now').isoformat() + "Z"
        })
    return vitals

def _generate_lab_panels(num=2):
    """Generates a list of random lab panel results."""
    panels = []
    lipid_panel_obs = {
        "High Density Lipoprotein Cholesterol": ("2085-9", "mg/dL", (30, 100)),
        "Low Density Lipoprotein Cholesterol": ("18262-6", "mg/dL", (50, 190)),
        "Total Cholesterol": ("2093-3", "mg/dL", (100, 300)),
        "Triglycerides": ("2571-8", "mg/dL", (50, 500))
    }
    cbc_panel_obs = {
        "Hemoglobin": ("718-7", "g/dL", (12, 18)),
        "WBC": ("6690-2", "10*3/uL", (4, 11)),
        "Platelets": ("777-3", "10*3/uL", (150, 450)),
        "Hematocrit": ("4544-3", "%", (35, 52))
    }

    for _ in range(num):
        reported_date = fake.date_time_between(start_date='-3y', end_date='now')
        panel_type = random.choice(["lipid", "cbc"])
        
        if panel_type == "lipid":
            panel_name, panel_code, observations_def = "Lipid Panel", "57698-3", lipid_panel_obs
        else:
            panel_name, panel_code, observations_def = "Complete blood count panel", "58410-2", cbc_panel_obs
            
        observations = []
        for obs_name, (obs_code, unit, val_range) in observations_def.items():
            observations.append({
                "codes": _generate_codes("loinc", obs_code, obs_name),
                "value": str(random.uniform(*val_range)),
                "unit": unit,
                "reportedTime": reported_date.isoformat() + "Z"
            })

        panels.append({
            "id": str(uuid.uuid4()),
            "reportedDate": reported_date.isoformat() + "Z",
            "category": "LAB",
            "codes": _generate_codes("loinc", panel_code, panel_name),
            "observations": observations
        })
    return panels

def _generate_conditions(num=1):
    """Generates a list of random conditions."""
    conditions = []
    cond_list = {
        "Hypertension": "59621000",
        "Type 2 diabetes mellitus": "44054006",
        "Hyperlipidemia": "55822004",
        "Polyp of colon": "68496003"
    }
    for _ in range(num):
        cond_name, cond_code = random.choice(list(cond_list.items()))
        conditions.append({
            "encounterId": str(uuid.uuid4()),
            "name": cond_name,
            "onsetDate": fake.date_time_between(start_date='-15y', end_date='-2y').isoformat() + "Z",
            "status": "active",
            "verificationStatus": "confirmed",
            "codes": _generate_codes("snomed.info/sct", cond_code, cond_name)
        })
    return conditions

def generate_insert_statement(table_name="public.t1d_observations"):
    """
    Generates a single SQL INSERT statement with complex, random patient data.
    """
    # Generate base UUIDs
    observation_id = str(uuid.uuid4())
    patient_id = str(uuid.uuid4())
    group_id = str(uuid.uuid4())
    mrn_id = "mrn-" + str(uuid.uuid4())
    cohort_id = str(uuid.uuid4())
    batch_id = str(uuid.uuid4())
    subcohort_id = str(uuid.uuid4())

    # Generate timestamps
    now = datetime.utcnow()
    created_time = fake.date_time_between(start_date='-1y', end_date='now')
    observed_time = created_time - timedelta(hours=random.randint(1, 24))
    last_modified_time = created_time + timedelta(minutes=random.randint(1, 60))

    # Create the complex 'patient' JSON object
    patient_data = {
        "id": patient_id,
        "ids": [{"id": mrn_id, "keyspace": "mrn"}],
        "mrn": mrn_id,
        "firstName": fake.first_name(),
        "lastName": fake.last_name(),
        "dob": fake.date_of_birth(minimum_age=30, maximum_age=90).strftime('%Y-%m-%dT00:00:00Z'),
        "dod": None,
        "gender": random.choice(["male", "female"]),
        "race": random.choice(["White", "Black or African American", "Asian", "Other"]),
        "ethnicity": random.choice(["Hispanic or Latino", "Not Hispanic or Latino"]),
        "maritalStatus": random.choice(["M", "S", "D", "W"]),
        "phoneNumber": "",
        "lastEmrPull": now.isoformat() + "Z",
        "medications": _generate_medications(num=random.randint(1, 3)),
        "vitals": _generate_vitals(num=random.randint(5, 8)),
        "labPanelResults": _generate_lab_panels(num=random.randint(1, 2)),
        "conditions": _generate_conditions(num=random.randint(100000, 200000)),
        "procedures": [], # Keeping this simple for now, can be expanded
        "encounters": [] # Keeping this simple for now, can be expanded
    }
    patient_json_string = json.dumps(patient_data).replace("'", "''")

    # Assemble the final SQL statement
    sql_statement = f"""
-- Inserting new row with complex patient data
INSERT INTO {table_name} (
    observation_id, ext_observation_id, group_id, patient_id, patient, observed, created, status, last_modified, last_modified_by,
    priority_level, priority_value, score, cohort_id, batch_id, pcp_name, subcohort_filter_id, location, explainability, model_info
) VALUES (
    '{observation_id}', 'ext_{observation_id[:8]}', '{group_id}', '{patient_id}',
    '{patient_json_string}'::jsonb,
    '{observed_time.isoformat()}Z', '{created_time.isoformat()}Z', {random.randint(0, 3)}, '{last_modified_time.isoformat()}Z', 'python_script_loader',
    {random.randint(0, 5)}, {random.randint(0, 10)}, {random.uniform(0.0, 1.0):.6f},
    '{cohort_id}', '{batch_id}', 'Dr. {fake.last_name()}', '{subcohort_id}',
    '{fake.city().replace("'", "''")}', '[]', '[]'
);"""
    return sql_statement

# --- Main execution ---
if __name__ == "__main__":
    number_of_rows = 1000
    
    print(f"-- Generated SQL script for {number_of_rows} row(s).")
    with open("out.sql", "w") as f: 
        f.writelines(generate_insert_statement() for _ in range(number_of_rows))

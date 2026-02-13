- having to pre-declare all data involved in the choreography is maybe not ideal? 
- think about how a toy version of reveal would look like a choreography...
Dataflow: 
Cohort --> PatientIngested --> DataReady --> WorkflowComplete --> ObservationIngested --> ObservationProcessed

App structure: 
Iris --|--> Kelpie --|--> Hummus--|--> Sullivan --|
^______|      ^______|     ^______|      ^________|


1. (IRIS) Get list of mrns from DATASOURCE
2. (IRIS) For each mrn, create patient id if not exists (HUMMUS) + emit data-ready to DATA_READY_QUEUE
3. (KELPIE) Pull from DATA_READY_QUEUE + start WORKFLOW based on DATASOURCE
4. (KELPIE) Field WORKFLOW_STEP commands by extracting static data (KELPIE), persisting data (KELPIE), transforming inputs (INPUT_ADAPTER), transforming outputs (OUTPUT_ADATPER), or running syncs (HUMMUS)
 - OUTPUT_ADATPER: respond to requests back to KELPIE
 - INPUT_ADATPER: respond to requests back to KELPIE
 - HUMMUS: submits sync jobs to itself 
5. 
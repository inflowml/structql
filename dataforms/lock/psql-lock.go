package lock

//AccessShare conflicts with the ACCESS EXCLUSIVE lock mode only.
const AccessShare = "ACCESS SHARE"

//RowShare conflicts with the EXCLUSIVE and ACCESS EXCLUSIVE lock modes.
const RowShare = "ROW SHARE"

//RowExclusive conflicts with the SHARE, SHARE ROW EXCLUSIVE, EXCLUSIVE, and ACCESS EXCLUSIVE lock modes.
const RowExclusive = "ROW EXCLUSIVE"

//ShareUpdate conflicts with the SHARE UPDATE EXCLUSIVE, SHARE,
//SHARE ROW EXCLUSIVE, EXCLUSIVE, and ACCESS EXCLUSIVE lock modes.
//This mode protects a table against concurrent schema changes and VACUUM runs.
const ShareUpdate = "SHARE UPDATE EXCLUSIVE"

//Share conflicts with the ROW EXCLUSIVE, SHARE UPDATE EXCLUSIVE,
//SHARE ROW EXCLUSIVE, EXCLUSIVE, and ACCESS EXCLUSIVE lock modes.
//This mode protects a table against concurrent data changes.
const Share = "SHARE"

//ShareRow conflicts with the ROW EXCLUSIVE, SHARE UPDATE EXCLUSIVE,
//SHARE, SHARE ROW EXCLUSIVE, EXCLUSIVE, and ACCESS EXCLUSIVE lock modes.
//This mode protects a table against concurrent data changes,
//and is self-exclusive so that only one session can hold it at a time.
const ShareRow = "SHARE ROW EXCLUSIVE"

//Exclusive conflicts with the ROW SHARE, ROW EXCLUSIVE, SHARE UPDATE EXCLUSIVE,
//SHARE, SHARE ROW EXCLUSIVE, EXCLUSIVE, and ACCESS EXCLUSIVE lock modes.
//This mode allows only concurrent ACCESS SHARE locks,
//i.e., only reads from the table can proceed in parallel with a transaction holding this lock mode.
const Exclusive = "EXCLUSIVE"

//AccessExclusive Conflicts with locks of all modes
//(ACCESS SHARE, ROW SHARE, ROW EXCLUSIVE, SHARE UPDATE EXCLUSIVE,
//SHARE, SHARE ROW EXCLUSIVE, EXCLUSIVE, and ACCESS EXCLUSIVE).
//This mode guarantees that the holder is the only transaction accessing the table in any way.
const AccessExclusive = "ACCESS EXCLUSIVE"

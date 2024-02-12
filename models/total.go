package models

type AttendanceTotal struct {
	All     int64        `json:"all"`
	Today   AttendanceWM `json:"today"`
	Weekly  AttendanceWM `json:"weekly"`
	Monthly AttendanceWM `json:"monthly"`
}

type AttendanceWM struct {
	Present int64 `json:"present"`
	Absent  int64 `json:"absent"`
}

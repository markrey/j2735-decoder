/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "AddGrpB"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#ifndef	_AddGrpB_TimeMark_H_
#define	_AddGrpB_TimeMark_H_


#include <asn_application.h>

/* Including external dependencies */
#include "Year.h"
#include "Month.h"
#include "Day.h"
#include "SummerTime.h"
#include "Holiday.h"
#include "DayOfWeek.h"
#include "Hour.h"
#include "Minute.h"
#include "Second.h"
#include "TenthSecond.h"
#include <constr_SEQUENCE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* TimeMark */
typedef struct AddGrpB_TimeMark {
	Year_t	 year;
	Month_t	 month;
	Day_t	 day;
	SummerTime_t	 summerTime;
	Holiday_t	 holiday;
	DayOfWeek_t	 dayofWeek;
	Hour_t	 hour;
	Minute_t	 minute;
	Second_t	 second;
	TenthSecond_t	 tenthSecond;
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} AddGrpB_TimeMark_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_AddGrpB_TimeMark;

#ifdef __cplusplus
}
#endif

#endif	/* _AddGrpB_TimeMark_H_ */
#include <asn_internal.h>

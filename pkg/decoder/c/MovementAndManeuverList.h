/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "DSRC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#ifndef	_MovementAndManeuverList_H_
#define	_MovementAndManeuverList_H_


#include <asn_application.h>

/* Including external dependencies */
#include <asn_SEQUENCE_OF.h>
#include <constr_SEQUENCE_OF.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Forward declarations */
struct MovementAndManeuver;

/* MovementAndManeuverList */
typedef struct MovementAndManeuverList {
	A_SEQUENCE_OF(struct MovementAndManeuver) list;
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} MovementAndManeuverList_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_MovementAndManeuverList;
extern asn_SET_OF_specifics_t asn_SPC_MovementAndManeuverList_specs_1;
extern asn_TYPE_member_t asn_MBR_MovementAndManeuverList_1[1];
extern asn_per_constraints_t asn_PER_type_MovementAndManeuverList_constr_1;

#ifdef __cplusplus
}
#endif

/* Referred external types */
#include "MovementAndManeuver.h"

#endif	/* _MovementAndManeuverList_H_ */
#include <asn_internal.h>

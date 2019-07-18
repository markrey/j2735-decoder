/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "DSRC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#ifndef	_PrivilegedEventFlags_H_
#define	_PrivilegedEventFlags_H_


#include <asn_application.h>

/* Including external dependencies */
#include <BIT_STRING.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Dependencies */
typedef enum PrivilegedEventFlags {
	PrivilegedEventFlags_peUnavailable	= 0,
	PrivilegedEventFlags_peEmergencyResponse	= 1,
	PrivilegedEventFlags_peEmergencyLightsActive	= 2,
	PrivilegedEventFlags_peEmergencySoundActive	= 3,
	PrivilegedEventFlags_peNonEmergencyLightsActive	= 4,
	PrivilegedEventFlags_peNonEmergencySoundActive	= 5
} e_PrivilegedEventFlags;

/* PrivilegedEventFlags */
typedef BIT_STRING_t	 PrivilegedEventFlags_t;

/* Implementation */
extern asn_per_constraints_t asn_PER_type_PrivilegedEventFlags_constr_1;
extern asn_TYPE_descriptor_t asn_DEF_PrivilegedEventFlags;
asn_struct_free_f PrivilegedEventFlags_free;
asn_struct_print_f PrivilegedEventFlags_print;
asn_constr_check_f PrivilegedEventFlags_constraint;
ber_type_decoder_f PrivilegedEventFlags_decode_ber;
der_type_encoder_f PrivilegedEventFlags_encode_der;
xer_type_decoder_f PrivilegedEventFlags_decode_xer;
xer_type_encoder_f PrivilegedEventFlags_encode_xer;
oer_type_decoder_f PrivilegedEventFlags_decode_oer;
oer_type_encoder_f PrivilegedEventFlags_encode_oer;
per_type_decoder_f PrivilegedEventFlags_decode_uper;
per_type_encoder_f PrivilegedEventFlags_encode_uper;

#ifdef __cplusplus
}
#endif

#endif	/* _PrivilegedEventFlags_H_ */
#include <asn_internal.h>

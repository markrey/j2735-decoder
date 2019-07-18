/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "DSRC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#include "VehicleClassification.h"

static int
memb_regional_constraint_1(const asn_TYPE_descriptor_t *td, const void *sptr,
			asn_app_constraint_failed_f *ctfailcb, void *app_key) {
	size_t size;
	
	if(!sptr) {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: value not given (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
	
	/* Determine the number of elements */
	size = _A_CSEQUENCE_FROM_VOID(sptr)->count;
	
	if((size >= 1 && size <= 4)) {
		/* Perform validation of the inner elements */
		return td->encoding_constraints.general_constraints(td, sptr, ctfailcb, app_key);
	} else {
		ASN__CTFAIL(app_key, td, sptr,
			"%s: constraint failed (%s:%d)",
			td->name, __FILE__, __LINE__);
		return -1;
	}
}

static asn_oer_constraints_t asn_OER_type_regional_constr_10 CC_NOTUSED = {
	{ 0, 0 },
	-1	/* (SIZE(1..4)) */};
static asn_per_constraints_t asn_PER_type_regional_constr_10 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 2,  2,  1,  4 }	/* (SIZE(1..4)) */,
	0, 0	/* No PER value map */
};
static asn_oer_constraints_t asn_OER_memb_regional_constr_10 CC_NOTUSED = {
	{ 0, 0 },
	-1	/* (SIZE(1..4)) */};
static asn_per_constraints_t asn_PER_memb_regional_constr_10 CC_NOTUSED = {
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	{ APC_CONSTRAINED,	 2,  2,  1,  4 }	/* (SIZE(1..4)) */,
	0, 0	/* No PER value map */
};
static asn_TYPE_member_t asn_MBR_regional_10[] = {
	{ ATF_POINTER, 0, 0,
		(ASN_TAG_CLASS_UNIVERSAL | (16 << 2)),
		0,
		&asn_DEF_RegionalExtension_116P0,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		""
		},
};
static const ber_tlv_tag_t asn_DEF_regional_tags_10[] = {
	(ASN_TAG_CLASS_CONTEXT | (8 << 2)),
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static asn_SET_OF_specifics_t asn_SPC_regional_specs_10 = {
	sizeof(struct VehicleClassification__regional),
	offsetof(struct VehicleClassification__regional, _asn_ctx),
	0,	/* XER encoding is XMLDelimitedItemList */
};
static /* Use -fall-defs-global to expose */
asn_TYPE_descriptor_t asn_DEF_regional_10 = {
	"regional",
	"regional",
	&asn_OP_SEQUENCE_OF,
	asn_DEF_regional_tags_10,
	sizeof(asn_DEF_regional_tags_10)
		/sizeof(asn_DEF_regional_tags_10[0]) - 1, /* 1 */
	asn_DEF_regional_tags_10,	/* Same as above */
	sizeof(asn_DEF_regional_tags_10)
		/sizeof(asn_DEF_regional_tags_10[0]), /* 2 */
	{ &asn_OER_type_regional_constr_10, &asn_PER_type_regional_constr_10, SEQUENCE_OF_constraint },
	asn_MBR_regional_10,
	1,	/* Single element */
	&asn_SPC_regional_specs_10	/* Additional specs */
};

asn_TYPE_member_t asn_MBR_VehicleClassification_1[] = {
	{ ATF_POINTER, 9, offsetof(struct VehicleClassification, keyType),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_BasicVehicleClass,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"keyType"
		},
	{ ATF_POINTER, 8, offsetof(struct VehicleClassification, role),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_BasicVehicleRole,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"role"
		},
	{ ATF_POINTER, 7, offsetof(struct VehicleClassification, iso3883),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_Iso3833VehicleType,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"iso3883"
		},
	{ ATF_POINTER, 6, offsetof(struct VehicleClassification, hpmsType),
		(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_VehicleType,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"hpmsType"
		},
	{ ATF_POINTER, 5, offsetof(struct VehicleClassification, vehicleType),
		(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_VehicleGroupAffected,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"vehicleType"
		},
	{ ATF_POINTER, 4, offsetof(struct VehicleClassification, responseEquip),
		(ASN_TAG_CLASS_CONTEXT | (5 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_IncidentResponseEquipment,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"responseEquip"
		},
	{ ATF_POINTER, 3, offsetof(struct VehicleClassification, responderType),
		(ASN_TAG_CLASS_CONTEXT | (6 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ResponderGroupAffected,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"responderType"
		},
	{ ATF_POINTER, 2, offsetof(struct VehicleClassification, fuelType),
		(ASN_TAG_CLASS_CONTEXT | (7 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_FuelType,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"fuelType"
		},
	{ ATF_POINTER, 1, offsetof(struct VehicleClassification, regional),
		(ASN_TAG_CLASS_CONTEXT | (8 << 2)),
		0,
		&asn_DEF_regional_10,
		0,
		{ &asn_OER_memb_regional_constr_10, &asn_PER_memb_regional_constr_10,  memb_regional_constraint_1 },
		0, 0, /* No default value */
		"regional"
		},
};
static const int asn_MAP_VehicleClassification_oms_1[] = { 0, 1, 2, 3, 4, 5, 6, 7, 8 };
static const ber_tlv_tag_t asn_DEF_VehicleClassification_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_VehicleClassification_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* keyType */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* role */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 }, /* iso3883 */
    { (ASN_TAG_CLASS_CONTEXT | (3 << 2)), 3, 0, 0 }, /* hpmsType */
    { (ASN_TAG_CLASS_CONTEXT | (4 << 2)), 4, 0, 0 }, /* vehicleType */
    { (ASN_TAG_CLASS_CONTEXT | (5 << 2)), 5, 0, 0 }, /* responseEquip */
    { (ASN_TAG_CLASS_CONTEXT | (6 << 2)), 6, 0, 0 }, /* responderType */
    { (ASN_TAG_CLASS_CONTEXT | (7 << 2)), 7, 0, 0 }, /* fuelType */
    { (ASN_TAG_CLASS_CONTEXT | (8 << 2)), 8, 0, 0 } /* regional */
};
asn_SEQUENCE_specifics_t asn_SPC_VehicleClassification_specs_1 = {
	sizeof(struct VehicleClassification),
	offsetof(struct VehicleClassification, _asn_ctx),
	asn_MAP_VehicleClassification_tag2el_1,
	9,	/* Count of tags in the map */
	asn_MAP_VehicleClassification_oms_1,	/* Optional members */
	9, 0,	/* Root/Additions */
	9,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_VehicleClassification = {
	"VehicleClassification",
	"VehicleClassification",
	&asn_OP_SEQUENCE,
	asn_DEF_VehicleClassification_tags_1,
	sizeof(asn_DEF_VehicleClassification_tags_1)
		/sizeof(asn_DEF_VehicleClassification_tags_1[0]), /* 1 */
	asn_DEF_VehicleClassification_tags_1,	/* Same as above */
	sizeof(asn_DEF_VehicleClassification_tags_1)
		/sizeof(asn_DEF_VehicleClassification_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_VehicleClassification_1,
	9,	/* Elements count */
	&asn_SPC_VehicleClassification_specs_1	/* Additional specs */
};


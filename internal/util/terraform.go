package util

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// UpdateTerraformObjectWithAttr adds an Attribute to a Terraform object
func UpdateTerraformObjectWithAttr(ctx context.Context, obj types.Object, key string, value attr.Value) (types.Object, diag.Diagnostics) {
	attrTypes := obj.AttributeTypes(ctx)
	attrValues := obj.Attributes()
	attrValues[key] = value
	return types.ObjectValue(attrTypes, attrValues)
}

func CompareTerraformAttrToString(attr attr.Value, s string) bool {
	if sAttr, ok := attr.(types.String); ok {
		return sAttr.ValueString() == s
	}
	return false
}

// MergeTerraformObjects combines two Terraform Objects replacing any null or unknown attribute values in `old` with
// matching attributes from `new`.  Object type attributes are handled recursively to avoid overriding existing
// nested attributes in the old Object. Full type information must be specified.
//
// The reason for this function is to handle situations where a remote resource was created but not all configuration
// was performed successfully.  Instead of deleting the misconfigured resource, we can warn the user, and allow them
// to fix the configuration.  In the case of Pulsar namespaces, it's possible that a namespace has been created, but
// not all of the policy configuration was completed successfully.  If the user is warned of the issues, they can
// re-sync their remote state, and then decide how to proceed, either changing the configuration or deleting the namespace.
func MergeTerraformObjects(old, new types.Object, attributeTypes map[string]attr.Type) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if attributeTypes == nil {
		diags.AddWarning("Failed to merge state objects", "No type information provided for object: "+old.String())
		return old, diags
	}
	if old.IsNull() || old.IsUnknown() {
		return basetypes.NewObjectValueMust(attributeTypes, new.Attributes()), diags
	}
	oldAttributes := old.Attributes()
	newAttributes := new.Attributes()
	attributes := map[string]attr.Value{}
	for name, newValue := range newAttributes {

		oldValue, ok := oldAttributes[name]
		if !ok || oldValue.IsNull() || oldValue.IsUnknown() {
			attributes[name] = newValue
			continue
		}

		if oldObjValue, ok := oldValue.(types.Object); ok {
			newObjValue, ok := newValue.(types.Object)
			if !ok {
				diags.AddWarning("Non matching types for attribute: "+name,
					fmt.Sprintf("Existing object attribute can't be replaced with type `%v`", newValue.Type(context.Background()).String()))
				continue
			}
			typeInfo, ok := attributeTypes[name].(types.ObjectType)
			if !ok {
				diags.AddWarning("Missing type information for attribute "+name, "No type information found when merging objects")
				continue
			}
			newObjValue, mergeDiags := MergeTerraformObjects(oldObjValue, newObjValue, typeInfo.AttributeTypes())
			diags.Append(mergeDiags...)
			if diags.HasError() {
				return old, diags
			}
			attributes[name] = newObjValue
			continue
		}

		attributes[name] = oldValue
	}

	return basetypes.NewObjectValue(attributeTypes, attributes)
}

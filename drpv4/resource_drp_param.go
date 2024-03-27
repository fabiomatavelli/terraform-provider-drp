package drpv4

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.com/rackn/provision/v4/api"
	"gitlab.com/rackn/provision/v4/models"
)

var _ resource.Resource = &ParamResource{}
var _ resource.ResourceWithImportState = &ParamResource{}

func NewParamResource() resource.Resource {
	return &ParamResource{}
}

// ParamResource provides the Terraform schema for the param resource
type ParamResource struct {
	session *api.Client
}

type ParamResourceModel struct {
	Id types.String `tfsdk:"id"`

	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Documentation types.String `tfsdk:"documentation"`
	Schema        types.Map    `tfsdk:"schema"`
	Secure        types.Bool   `tfsdk:"secure"`
}

func (r *ParamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_param"
}

func (r *ParamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Param resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Param ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Param name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Param description",
				Optional:            true,
			},
			"documentation": schema.StringAttribute{
				MarkdownDescription: "Param documentation",
				Optional:            true,
			},
			"schema": schema.MapAttribute{
				MarkdownDescription: "Param schema",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"secure": schema.BoolAttribute{
				MarkdownDescription: "Param secure",
				Optional:            true,
			},
		},
	}
}

func (r *ParamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.session = resourceGenericConfigure(ctx, req, resp)
}

func (r *ParamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[resourceParam] Creating param")
	var plan ParamResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := &models.Param{
		Name:          plan.Name.ValueString(),
		Description:   plan.Description.ValueString(),
		Documentation: plan.Documentation.ValueString(),
		Schema:        plan.Schema.Elements(),
		Secure:        plan.Secure.ValueBool(),
	}

	if err := r.session.CreateModel(param); err != nil {
		tflog.Error(ctx, fmt.Sprintf("[resourceParamCreate] error creating param: %s", err))
		resp.Diagnostics.AddError(fmt.Sprintf("error creating param %s", param.Name), err.Error())
		return
	}

	plan.Id = types.StringValue(param.Key())
	resp.State.Set(ctx, plan)
}

func (r *ParamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[resourceParam] Reading param")
	var plan ParamResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	po, err := r.session.GetModel("params", plan.Id.ValueString())
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("[resourceParamRead] error reading param: %s", err))
		resp.Diagnostics.AddError(fmt.Sprintf("error reading param %s", plan.Id), err.Error())
		return
	}

	param := po.(*models.Param)
	plan.Name = types.StringValue(param.Name)
	plan.Description = types.StringValue(param.Description)
	plan.Documentation = types.StringValue(param.Documentation)
	plan.Schema, _ = types.MapValue(types.StringType, param.Schema.(map[string]attr.Value))
	plan.Secure = types.BoolValue(param.Secure)

	resp.State.Set(ctx, plan)
}

func (r *ParamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[resourceParam] Updating param")
	var plan ParamResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := &models.Param{
		Name:          plan.Name.ValueString(),
		Description:   plan.Description.ValueString(),
		Documentation: plan.Documentation.ValueString(),
		Schema:        plan.Schema.Elements(),
		Secure:        plan.Secure.ValueBool(),
	}

	if err := r.session.PutModel(param); err != nil {
		tflog.Error(ctx, fmt.Sprintf("[resourceParamUpdate] error updating param: %s", err))
		resp.Diagnostics.AddError(fmt.Sprintf("error updating param %s", param.Name), err.Error())
		return
	}

	resp.State.Set(ctx, plan)
}

func (r *ParamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[resourceParam] Deleting param")
	var plan ParamResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.session.DeleteModel("params", plan.Id.ValueString()); err != nil {
		tflog.Error(ctx, fmt.Sprintf("[resourceParamDelete] error deleting param: %s", err))
		resp.Diagnostics.AddError(fmt.Sprintf("error deleting param %s", plan.Id), err.Error())
		return
	}

	resp.State.Set(ctx, plan)
}

func (r *ParamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package openapi_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"dazzle/internal/domain"
	"dazzle/internal/infrastructure/openapi"
)

func fixturesDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "testdata", "fixtures")
}

func TestRepository_Load(t *testing.T) {
	repo := openapi.NewRepository()
	spec, err := repo.Load(context.Background(), filepath.Join(fixturesDir(), "petstore.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("spec info", func(t *testing.T) {
		if spec.Info.Title != "Petstore API" {
			t.Errorf("expected title 'Petstore API', got %q", spec.Info.Title)
		}
		if spec.Info.Version != "1.0.0" {
			t.Errorf("expected version '1.0.0', got %q", spec.Info.Version)
		}
	})

	t.Run("servers", func(t *testing.T) {
		if len(spec.Servers) != 1 {
			t.Fatalf("expected 1 server, got %d", len(spec.Servers))
		}
		if spec.Servers[0].URL != "https://api.petstore.example" {
			t.Errorf("unexpected server URL: %s", spec.Servers[0].URL)
		}
	})

	t.Run("operations count", func(t *testing.T) {
		if len(spec.Operations) != 4 {
			t.Errorf("expected 4 operations, got %d", len(spec.Operations))
		}
	})

	t.Run("operation fields", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		listPets, ok := ops["listPets"]
		if !ok {
			t.Fatal("missing listPets operation")
		}
		if listPets.Method != domain.GET {
			t.Errorf("expected GET, got %s", listPets.Method)
		}
		if listPets.Path != "/pets" {
			t.Errorf("expected /pets, got %s", listPets.Path)
		}
		if listPets.Summary != "List all pets" {
			t.Errorf("unexpected summary: %s", listPets.Summary)
		}
	})

	t.Run("operation-level parameters", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		listPets := ops["listPets"]
		if len(listPets.Parameters) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(listPets.Parameters))
		}
		param := listPets.Parameters[0]
		if param.Name != "limit" {
			t.Errorf("expected param name 'limit', got %q", param.Name)
		}
		if param.In != domain.ParameterInQuery {
			t.Errorf("expected param in 'query', got %q", param.In)
		}
	})

	t.Run("path-level parameters merged into operations", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		// getPet has no operation-level params but should inherit petId from path level
		getPet := ops["getPet"]
		if len(getPet.Parameters) != 1 {
			t.Fatalf("expected 1 parameter (from path level), got %d", len(getPet.Parameters))
		}
		if getPet.Parameters[0].Name != "petId" {
			t.Errorf("expected param name 'petId', got %q", getPet.Parameters[0].Name)
		}
		if getPet.Parameters[0].In != domain.ParameterInPath {
			t.Errorf("expected param in 'path', got %q", getPet.Parameters[0].In)
		}

		// deletePet should also inherit petId from path level
		deletePet := ops["deletePet"]
		if len(deletePet.Parameters) != 1 {
			t.Fatalf("expected 1 parameter (from path level), got %d", len(deletePet.Parameters))
		}
		if deletePet.Parameters[0].Name != "petId" {
			t.Errorf("expected param name 'petId', got %q", deletePet.Parameters[0].Name)
		}
	})

	t.Run("parameter schemas", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		limit := ops["listPets"].Parameters[0]
		if limit.Schema == nil {
			t.Fatal("expected limit param to have schema")
		}
		if limit.Schema.Type != domain.SchemaTypeInteger {
			t.Errorf("expected integer type, got %q", limit.Schema.Type)
		}
		if limit.Schema.Format != "int32" {
			t.Errorf("expected int32 format, got %q", limit.Schema.Format)
		}

		petId := ops["getPet"].Parameters[0]
		if petId.Schema == nil {
			t.Fatal("expected petId param to have schema")
		}
		if petId.Schema.Type != domain.SchemaTypeInteger {
			t.Errorf("expected integer type, got %q", petId.Schema.Type)
		}
	})

	t.Run("request body", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		createPet := ops["createPet"]
		if createPet.RequestBody == nil {
			t.Fatal("expected createPet to have a request body")
		}
		if !createPet.RequestBody.Required {
			t.Error("expected request body to be required")
		}
		jsonContent, ok := createPet.RequestBody.Content["application/json"]
		if !ok {
			t.Fatal("expected application/json content type")
		}
		if jsonContent.Schema == nil {
			t.Fatal("expected schema in json content")
		}
		if jsonContent.Schema.Type != domain.SchemaTypeObject {
			t.Errorf("expected object type, got %q", jsonContent.Schema.Type)
		}
		nameProp, ok := jsonContent.Schema.Properties["name"]
		if !ok {
			t.Fatal("expected 'name' property in schema")
		}
		if nameProp.Type != domain.SchemaTypeString {
			t.Errorf("expected string type for name, got %q", nameProp.Type)
		}
	})

	t.Run("responses", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		listPets := ops["listPets"]
		resp200, ok := listPets.Responses["200"]
		if !ok {
			t.Fatal("expected 200 response for listPets")
		}
		if resp200.Description != "A list of pets" {
			t.Errorf("unexpected description: %q", resp200.Description)
		}
		jsonContent, ok := resp200.Content["application/json"]
		if !ok {
			t.Fatal("expected application/json content in 200 response")
		}
		if jsonContent.Schema == nil {
			t.Fatal("expected schema in response content")
		}
		if jsonContent.Schema.Type != domain.SchemaTypeArray {
			t.Errorf("expected array type, got %q", jsonContent.Schema.Type)
		}
		if jsonContent.Schema.Items == nil {
			t.Fatal("expected items schema for array")
		}
		if jsonContent.Schema.Items.Type != domain.SchemaTypeObject {
			t.Errorf("expected object item type, got %q", jsonContent.Schema.Items.Type)
		}

		deletePet := ops["deletePet"]
		resp204, ok := deletePet.Responses["204"]
		if !ok {
			t.Fatal("expected 204 response for deletePet")
		}
		if resp204.Description != "Pet deleted" {
			t.Errorf("unexpected description: %q", resp204.Description)
		}
		if resp204.Content != nil {
			t.Error("expected no content for 204 response")
		}
	})

	t.Run("response headers", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		resp200 := ops["listPets"].Responses["200"]
		if len(resp200.Headers) != 1 {
			t.Fatalf("expected 1 response header, got %d", len(resp200.Headers))
		}
		h, ok := resp200.Headers["X-Total-Count"]
		if !ok {
			t.Fatal("expected X-Total-Count header")
		}
		if h.Description != "Total number of pets" {
			t.Errorf("unexpected header description: %q", h.Description)
		}
		if h.Schema == nil {
			t.Fatal("expected header to have schema")
		}
		if h.Schema.Type != domain.SchemaTypeInteger {
			t.Errorf("expected integer type, got %q", h.Schema.Type)
		}
	})

	t.Run("array item properties preserved", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		resp200 := ops["listPets"].Responses["200"]
		jsonContent := resp200.Content["application/json"]
		items := jsonContent.Schema.Items
		if len(items.Properties) != 2 {
			t.Fatalf("expected 2 item properties (id, name), got %d", len(items.Properties))
		}
		if items.Properties["id"].Type != domain.SchemaTypeInteger {
			t.Errorf("expected integer type for id, got %q", items.Properties["id"].Type)
		}
		if items.Properties["name"].Type != domain.SchemaTypeString {
			t.Errorf("expected string type for name, got %q", items.Properties["name"].Type)
		}
	})

	t.Run("nil request body for GET operations", func(t *testing.T) {
		ops := indexByID(spec.Operations)

		if ops["listPets"].RequestBody != nil {
			t.Error("expected nil request body for listPets")
		}
		if ops["getPet"].RequestBody != nil {
			t.Error("expected nil request body for getPet")
		}
	})
}

func TestRepository_Load_InvalidFile(t *testing.T) {
	repo := openapi.NewRepository()
	_, err := repo.Load(context.Background(), "nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func indexByID(ops []domain.Operation) map[string]domain.Operation {
	m := make(map[string]domain.Operation, len(ops))
	for _, op := range ops {
		m[op.ID] = op
	}
	return m
}

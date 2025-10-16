package hcl

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

var (
	ErrCouldntParseFile                = errors.New("couldn't parse file")
	ErrUnexpectedBodyType              = errors.New("unexpected body type")
	ErrModuleMissingLabel              = errors.New("module block missing label")
	ErrTemplateWithInterpolation       = errors.New("template expressions with interpolation are not supported")
	ErrUnsupportedExpressionType       = errors.New("unsupported expression type")
	ErrNullValueCannotBeConvertedToStr = errors.New("null values cannot be converted to string")
)

type TFModule struct {
	Name      string
	Attribute string
}

func ParseModules(path, attributeKey string, valueRegex *regexp.Regexp) ([]TFModule, error) {
	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("%w (%q): %s", ErrCouldntParseFile, path, diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, ErrUnexpectedBodyType
	}

	var modules []TFModule

	for _, block := range body.Blocks {
		if block.Type == "module" {
			if len(block.Labels) == 0 {
				return nil, fmt.Errorf("%w at %s", ErrModuleMissingLabel, block.DefRange())
			}
			moduleName := block.Labels[0]

			if attr, exists := block.Body.Attributes[attributeKey]; exists {
				attribute, err := extractStringValue(attr.Expr)
				if err != nil {
					return nil, fmt.Errorf("couldn't extract %s from module %q: %w", attributeKey, moduleName, err)
				}
				attribute = extractValue(attribute, valueRegex)
				modules = append(modules, TFModule{
					Name:      moduleName,
					Attribute: attribute,
				})
			}
		}
	}

	return modules, nil
}

func extractStringValue(expr hclsyntax.Expression) (string, error) {
	switch e := expr.(type) {
	case *hclsyntax.TemplateExpr:
		if len(e.Parts) == 1 {
			if lit, ok := e.Parts[0].(*hclsyntax.LiteralValueExpr); ok {
				return convertToString(lit.Val)
			}
		}
		return "", ErrTemplateWithInterpolation
	case *hclsyntax.LiteralValueExpr:
		return convertToString(e.Val)
	default:
		return "", fmt.Errorf("%w: %T", ErrUnsupportedExpressionType, expr)
	}
}

func convertToString(val cty.Value) (string, error) {
	if val.IsNull() {
		return "", ErrNullValueCannotBeConvertedToStr
	}

	if val.Type() == cty.String {
		return val.AsString(), nil
	}

	if val.Type() == cty.Bool {
		if val.True() {
			return "true", nil
		}
		return "false", nil
	}

	if val.Type() == cty.Number {
		bf := val.AsBigFloat()
		if bf.IsInt() {
			i, _ := bf.Int64()
			return fmt.Sprintf("%d", i), nil
		}
		f, _ := bf.Float64()
		return fmt.Sprintf("%g", f), nil
	}

	return val.GoString(), nil
}

func extractValue(value string, valueRegex *regexp.Regexp) string {
	if valueRegex == nil {
		return value
	}

	matches := valueRegex.FindStringSubmatch(value)
	if len(matches) > 1 {
		return matches[1]
	}

	return value
}

package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/robertjanetzko/LegendsBrowser2/backend/util"
)

func (e *Entity) Position(id int) *EntityPosition {
	for _, p := range e.EntityPosition {
		if p.Id_ == id {
			return p
		}
	}
	return &EntityPosition{Name_: "UNKNOWN POSITION"}
}

func (p *EntityPosition) GenderName(hf *HistoricalFigure) string {
	if hf.Female() && p.NameFemale != "" {
		return p.NameFemale
	} else if hf.Male() && p.NameMale != "" {
		return p.NameMale
	} else {
		return p.Name_
	}
}

func (hf *HistoricalFigure) Female() bool {
	return hf.Sex == 0 || hf.Caste == "FEMALE"
}

func (hf *HistoricalFigure) Male() bool {
	return hf.Sex == 1 || hf.Caste == "MALE"
}

type HistoricalEventDetails interface {
	RelatedToEntity(int) bool
	RelatedToHf(int) bool
	Html() string
	Type() string
}

type HistoricalEventCollectionDetails interface {
}

func containsInt(list []int, id int) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}
	return false
}

var world *DfWorld

func andList(list []string) string {
	if len(list) > 1 {
		return strings.Join(list[:len(list)-1], ", ") + " and " + list[len(list)-1]
	}
	return strings.Join(list, ", ")
}

func articled(s string) string {
	if ok, _ := regexp.MatchString("^([aeio]|un|ul).*", s); ok {
		return "an " + s
	}
	return "a " + s
}

func artifact(id int) string {
	if x, ok := world.Artifacts[id]; ok {
		return fmt.Sprintf(`<a class="artifact" href="/artifact/%d">%s</a>`, x.Id(), util.Title(x.Name()))
	}
	return "UNKNOWN ARTIFACT"
}

func entity(id int) string {
	if x, ok := world.Entities[id]; ok {
		return fmt.Sprintf(`<a class="entity" href="/entity/%d">%s</a>`, x.Id(), util.Title(x.Name()))
	}
	return "UNKNOWN ENTITY"
}

func hf(id int) string {
	if x, ok := world.HistoricalFigures[id]; ok {
		return fmt.Sprintf(`the %s <a class="hf" href="/hf/%d">%s</a>`, x.Race, x.Id(), util.Title(x.Name()))
	}
	return "UNKNOWN HISTORICAL FIGURE"
}

func hfList(ids []int) string {
	return andList(util.Map(ids, hf))
}

func pronoun(id int) string {
	if x, ok := world.HistoricalFigures[id]; ok {
		if x.Female() {
			return "she"
		}
	}
	return "he"
}

func site(id int, prefix string) string {
	if x, ok := world.Sites[id]; ok {
		return fmt.Sprintf(`%s <a class="site" href="/site/%d">%s</a>`, prefix, x.Id(), util.Title(x.Name()))
	}
	return "UNKNOWN SITE"
}

func structure(siteId, structureId int) string {
	if x, ok := world.Sites[siteId]; ok {
		if y, ok := x.Structures[structureId]; ok {
			return fmt.Sprintf(`<a class="structure" href="/site/%d/structure/%d">%s</a>`, siteId, structureId, util.Title(y.Name()))
		}
	}
	return "UNKNOWN STRUCTURE"
}

func property(siteId, propertyId int) string {
	if x, ok := world.Sites[siteId]; ok {
		if y, ok := x.SiteProperties[propertyId]; ok {
			if y.StructureId != -1 {
				return structure(siteId, y.StructureId)
			}
			return articled(y.Type.String())
		}
	}
	return "UNKNOWN PROPERTY"
}

func region(id int) string {
	if x, ok := world.Regions[id]; ok {
		return fmt.Sprintf(`<a class="region" href="/region/%d">%s</a>`, x.Id(), util.Title(x.Name()))
	}
	return "UNKNOWN REGION"
}

func location(siteId int, sitePrefix string, regionId int, regionPrefix string) string {
	if siteId != -1 {
		return site(siteId, sitePrefix)
	}
	if regionId != -1 {
		return regionPrefix + " " + region(regionId)
	}
	return ""
}

func identity(id int) string {
	if x, ok := world.Identities[id]; ok {
		return fmt.Sprintf(`<a class="identity" href="/region/%d">%s</a>`, x.Id(), util.Title(x.Name()))
	}
	return "UNKNOWN IDENTITY"
}

func feature(x *Feature) string {
	switch x.Type {
	case FeatureType_DancePerformance:
		return "a perfomance of " + danceForm(x.Reference)
	case FeatureType_Images:
		if x.Reference != -1 {
			return "images of " + hf(x.Reference)
		}
		return "images"
	case FeatureType_MusicalPerformance:
		return "a perfomance of " + musicalForm(x.Reference)
	case FeatureType_PoetryRecital:
		return "a recital of " + poeticForm(x.Reference)
	case FeatureType_Storytelling:
		if x.Reference != -1 {
			return "a telling of the story of " + hf(x.Reference)
		}
		return "a story recital"
	default:
		return strcase.ToDelimited(x.Type.String(), ' ')
	}
}

func danceForm(id int) string {
	if x, ok := world.DanceForms[id]; ok {
		return fmt.Sprintf(`<a class="artform" href="/danceForm/%d">%s</a>`, id, util.Title(x.Name()))
	}
	return "UNKNOWN DANCE FORM"
}

func musicalForm(id int) string {
	if x, ok := world.MusicalForms[id]; ok {
		return fmt.Sprintf(`<a class="artform" href="/musicalForm/%d">%s</a>`, id, util.Title(x.Name()))
	}
	return "UNKNOWN MUSICAL FORM"
}

func poeticForm(id int) string {
	if x, ok := world.PoeticForms[id]; ok {
		return fmt.Sprintf(`<a class="artform" href="/poeticForm/%d">%s</a>`, id, util.Title(x.Name()))
	}
	return "UNKNOWN POETIC FORM"
}

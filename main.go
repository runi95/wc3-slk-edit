package main

import (
	"fmt"
	"github.com/andlabs/ui"
	"github.com/runi95/wts-parser/models"
	"github.com/runi95/wts-parser/parser"
	"gopkg.in/volatiletech/null.v6"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var baseUnitMap map[string]*models.SLKUnit
var unitFuncMap map[string]*models.UnitFunc
var selectedUnit null.Int
var lastValidIndex int

var mainWindow *ui.Window
var mh = new(modelHandler)
var uForm uiForm
var wFormOne weaponForm
var wFormTwo weaponForm
var dForm dataForm
var oForm otherForm
var races = []string{"\"_\"", "\"commoner\"", "\"creeps\"", "\"critters\"", "\"demon\"", "\"human\"", "\"naga\"", "\"nightelf\"", "\"orc\"", "\"other\"", "\"unknown\"", "\"undead\""}
var moveTypes = []string{"\"_\"", "\"foot\"", "\"horse\"", "\"fly\"", "\"hover\"", "\"float\"", "\"amph\""}
var pathingTextures = []string{"\"PathTextures\\10x10Simple.tga\"", "\"PathTextures\\10x10Simple.tga\"", "\"PathTextures\\10x10Simple.tga\"", "\"PathTextures\\10x10Simple.tga\"", "\"PathTextures\\12x12Simple.tga\"", "\"PathTextures\\6x6SimpleSolid.tga\"", "\"PathTextures\\4x4SimpleSolid.tga\"", "\"PathTextures\\16x16Simple.tga\"", "\"PathTextures\\DemonGatePath.tga\"", "\"PathTextures\\DarkPortalSE.tga\"", "\"PathTextures\\DarkPortalSW.tga\"", "\"PathTextures\\16x16Goldmine.tga\"", "\"PathTextures\\16x16Simple.tga\"", "\"PathTextures\\UndeadNecropolis.tga\"", "\"PathTextures\\16x16Simple.tga\"", "\"PathTextures\\6x6SimpleSolid.tga\"", "\"PathTextures\\8x8SimpleSolid.tga\"", "\"PathTextures\\12x12Simple.tga\"", "\"PathTextures\\6x6SimpleSolid.tga\"", "\"PathTextures\\16x16Simple.tga\""}
var weaponTypes = []string{"\"_\"", "\"normal\"", "\"instant\"", "\"artillery\"", "\"aline\"", "\"missile\"", "\"msplash\"", "\"mbounce\"", "\"mline\""}
var attackTypes = []string{"\"_\"", "\"normal\"", "\"pierce\"", "\"siege\"", "\"spells\"", "\"chaos\"", "\"magic\"", "\"hero\""}

// const MAXINT = 9999999

type modelHandler struct {
	rows          int
	slkUnitIdList []string
}

type targetGrid struct {
	*ui.Grid
	air          *ui.Checkbox
	alive        *ui.Checkbox
	allies       *ui.Checkbox
	ancient      *ui.Checkbox
	bridge       *ui.Checkbox
	dead         *ui.Checkbox
	debris       *ui.Checkbox
	decoration   *ui.Checkbox
	enemies      *ui.Checkbox
	friend       *ui.Checkbox
	ground       *ui.Checkbox
	hero         *ui.Checkbox
	invulnerable *ui.Checkbox
	item         *ui.Checkbox
	mechanical   *ui.Checkbox
	neutral      *ui.Checkbox
	nonancient   *ui.Checkbox
	nonhero      *ui.Checkbox
	nonsapper    *ui.Checkbox
	none         *ui.Checkbox
	notself      *ui.Checkbox
	organic      *ui.Checkbox
	player       *ui.Checkbox
	self         *ui.Checkbox
	structure    *ui.Checkbox
	sapper       *ui.Checkbox
	terrain      *ui.Checkbox
	tree         *ui.Checkbox
	vulnerable   *ui.Checkbox
	wall         *ui.Checkbox
	ward         *ui.Checkbox
}

type uiForm struct {
	unitId             *ui.Entry
	name               *ui.Entry
	abilities          *ui.Entry
	icon               *ui.Entry
	buttonX            *ui.Entry
	buttonY            *ui.Entry
	model              *ui.Entry
	soundSet           *ui.Entry
	hideMinimapDisplay *ui.Checkbox
	scalingValue       *ui.Entry
	selectionScale     *ui.Entry
	pathingTexture     *ui.Combobox
	color              *ui.ColorButton
	red                *ui.Slider
	green              *ui.Slider
	blue               *ui.Slider
	hotkey             *ui.Entry
	tooltip            *ui.Entry
	description        *ui.MultilineEntry
}

type weaponForm struct {
	enableWeapon     *ui.Checkbox
	backswingPoint   *ui.Entry
	damagePoint      *ui.Entry
	attackType       *ui.Combobox
	targets          targetGrid
	cooldown         *ui.Entry
	damageBase       *ui.Entry
	damageDice       *ui.Entry
	damageSides      *ui.Entry
	weaponRange      *ui.Entry
	weaponType       *ui.Combobox
	aoeFull          *ui.Entry
	aoeMedium        *ui.Entry
	aoeSmall         *ui.Entry
	aoeTargets       targetGrid
	projectile       *ui.Entry
	projectileHoming *ui.Checkbox
	projectileSpeed  *ui.Entry
}

type weaponTab struct {
	*ui.Tab
	weaponOne weaponForm
	weaponTwo weaponForm
}

type dataForm struct {
	acquisition          *ui.Entry
	builds               *ui.Entry
	sells                *ui.Entry
	upgradesTo           *ui.Entry
	trains               *ui.Entry
	health               *ui.Entry
	mana                 *ui.Entry
	isBuilding           *ui.Checkbox
	defense              *ui.Entry
	defenseType          *ui.Combobox
	lumberCost           *ui.Entry
	goldCost             *ui.Entry
	points               *ui.Entry
	race                 *ui.Combobox
	foodCost             *ui.Entry
	foodProduction       *ui.Entry
	isCampaign           *ui.Checkbox
	movementType         *ui.Combobox
	movementSpeed        *ui.Entry
	movementSpeedMinimum *ui.Entry
	movementSpeedMaximum *ui.Entry
	flyingHeight         *ui.Entry
	minimumFlyingHeight  *ui.Entry
}

type otherForm struct {
	deathType             *ui.Entry
	death                 *ui.Entry
	cargoSize             *ui.Entry
	turnRate              *ui.Entry
	canSleep              *ui.Checkbox
	canBeBuiltOn          *ui.Checkbox
	canBuildOn            *ui.Checkbox
	dropsItems            *ui.Checkbox
	elevationSamplePoints *ui.Entry
	elevationSampleRadius *ui.Entry
	targetedAs            targetGrid
	level                 *ui.Entry
}

func main() {
	log.Println("Reading UnitAbilities.slk...")

	unitAbilitiesBytes, err := ioutil.ReadFile("input/UnitAbilities.slk")
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}

	unitAbilitiesMap := parser.SlkToUnitAbilities(unitAbilitiesBytes)

	log.Println("Reading UnitData.slk...")

	unitDataBytes, err := ioutil.ReadFile("input/UnitData.slk")
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}

	unitDataMap := parser.SlkToUnitData(unitDataBytes)

	log.Println("Reading UnitUI.slk...")

	unitUIBytes, err := ioutil.ReadFile("input/UnitUI.slk")
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}

	unitUIMap := parser.SLKToUnitUI(unitUIBytes)

	log.Println("Reading UnitWeapons.slk...")

	unitWeaponsBytes, err := ioutil.ReadFile("input/UnitWeapons.slk")
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}

	unitWeaponsMap := parser.SLKToUnitWeapons(unitWeaponsBytes)

	log.Println("Reading UnitBalance.slk...")

	unitBalanceBytes, err := ioutil.ReadFile("input/UnitBalance.slk")
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}

	unitBalanceMap := parser.SLKToUnitBalance(unitBalanceBytes)

	log.Println("Reading CampaignUnitFunc.txt...")

	campaignUnitFuncBytes, err := ioutil.ReadFile("input/CampaignUnitFunc.txt")
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}

	unitFuncMap = parser.TxtToUnitFunc(campaignUnitFuncBytes)

	baseUnitMap = make(map[string]*models.SLKUnit)
	mh.slkUnitIdList = make([]string, len(unitDataMap))
	i := 0
	for k := range unitDataMap {
		slkUnit := new(models.SLKUnit)
		slkUnit.UnitAbilities = unitAbilitiesMap[k]
		slkUnit.UnitData = unitDataMap[k]
		slkUnit.UnitUI = unitUIMap[k]
		slkUnit.UnitWeapons = unitWeaponsMap[k]
		slkUnit.UnitBalance = unitBalanceMap[k]

		unitId := strings.Replace(k, "\"", "", -1)
		baseUnitMap[unitId] = slkUnit
		mh.slkUnitIdList[i] = unitId
		i++
	}
	mh.rows = len(mh.slkUnitIdList)

	ui.Main(setupUI)
}

func intToHex(i int) string {
	if i < 10 {
		return fmt.Sprint(i)
	} else if i < 16 {
		return fmt.Sprint(string(55 + i))
	} else {
		return ""
	}
}

func getNextValidUnitId(offset int) string {
	var str string

	if offset > 16383 {
		log.Println("Ran out of valid generated unit id's")
		return ""
	}

	var firstChar string

	switch offset / 4096 {
	case 0:
		firstChar = "u"
	case 1:
		firstChar = "n"
	case 2:
		firstChar = "h"
	case 3:
		firstChar = "o"
	case 4:
		firstChar = "e"
	default:
		firstChar = "u"
	}

	str = firstChar + intToHex(offset/256) + intToHex(int(offset/16)%16) + intToHex(offset%16)
	if _, ok := unitFuncMap[str]; !ok {
		lastValidIndex = offset
		return str
	}

	return getNextValidUnitId(offset + 1)
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""), // column 0 text
		ui.TableColor{},    // row background color
	}
}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return mh.rows
}

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	return ui.TableString(mh.slkUnitIdList[row] + " - " + unitFuncMap[mh.slkUnitIdList[row]].Name.String)
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	selectedUnit.SetValid(row)
	unit := baseUnitMap[mh.slkUnitIdList[row]]
	unitFunc := unitFuncMap[mh.slkUnitIdList[row]]

	// Set UI Form
	uForm.name.SetText(unitFunc.Name.String)
	uForm.unitId.SetText(unitFunc.UnitId)
	desc := strings.Replace(unitFunc.Ubertip.String, "|n", "\n", -1)
	desc = strings.Replace(desc, "\"", "", -1)
	uForm.description.SetText(desc)
	uForm.hotkey.SetText(unitFunc.Hotkey.String)
	uForm.icon.SetText(unitFunc.Art.String)
	uForm.tooltip.SetText(unitFunc.Tip.String)
	buttonSplit := strings.Split(unitFunc.Buttonpos.String, ",")
	if len(buttonSplit) > 1 {
		uForm.buttonX.SetText(buttonSplit[0])
		uForm.buttonY.SetText(buttonSplit[1])
	}
	abilities := strings.Replace(unit.UnitAbilities.AbilList.String, "\"", "", -1)
	uForm.abilities.SetText(abilities)
	if unit.UnitUI.HideOnMinimap.String == "1" {
		uForm.hideMinimapDisplay.SetChecked(true)
	} else {
		uForm.hideMinimapDisplay.SetChecked(false)
	}
	red, err := strconv.Atoi(unit.UnitUI.Red.String)
	if err != nil {
		log.Println(err)
	}
	green, err := strconv.Atoi(unit.UnitUI.Green.String)
	if err != nil {
		log.Println(err)
	}
	blue, err := strconv.Atoi(unit.UnitUI.Blue.String)
	if err != nil {
		log.Println(err)
	}
	uForm.red.SetValue(red)
	uForm.green.SetValue(green)
	uForm.blue.SetValue(blue)
	uForm.color.SetColor(float64(uForm.red.Value())*0.003921569, float64(uForm.green.Value())*0.003921569, float64(uForm.blue.Value())*0.003921569, 1)
	model := strings.Replace(unit.UnitUI.File.String, "\"", "", -1)
	uForm.model.SetText(model)
	uForm.scalingValue.SetText(unit.UnitUI.ModelScale.String)
	uForm.selectionScale.SetText(unit.UnitUI.Scale.String)
	soundSet := strings.Replace(unit.UnitUI.UnitSound.String, "\"", "", -1)
	uForm.soundSet.SetText(soundSet)
	pathTex := strings.Replace(unit.UnitData.PathTex.String, "\"", "", -1)
	var pathingTextureSelected null.Int
	switch pathTex {
	case "PathTextures\\10x10Simple.tga":
		pathingTextureSelected.SetValid(0)
	case "PathTextures\\12x12Simple.tga":
		pathingTextureSelected.SetValid(4)
	case "PathTextures\\6x6SimpleSolid.tga":
		pathingTextureSelected.SetValid(5)
	case "PathTextures\\4x4SimpleSolid.tga":
		pathingTextureSelected.SetValid(6)
	case "PathTextures\\16x16Simple.tga":
		pathingTextureSelected.SetValid(7)
	case "PathTextures\\DemonGatePath.tga":
		pathingTextureSelected.SetValid(8)
	case "PathTextures\\DarkPortalSE.tga":
		pathingTextureSelected.SetValid(9)
	case "PathTextures\\DarkPortalSW.tga":
		pathingTextureSelected.SetValid(10)
	case "PathTextures\\16x16Goldmine.tga":
		pathingTextureSelected.SetValid(11)
	case "PathTextures\\UndeadNecropolis.tga":
		pathingTextureSelected.SetValid(13)
	case "PathTextures\\8x8SimpleSolid.tga":
		pathingTextureSelected.SetValid(16)
	default:
		pathingTextureSelected.SetValid(6)
	}
	if pathingTextureSelected.Valid {
		uForm.pathingTexture.SetSelected(pathingTextureSelected.Int)
	}

	// Set Weapon Form 1
	if unit.UnitWeapons.WeapsOn.Valid && unit.UnitWeapons.WeapsOn.String == "1" || unit.UnitWeapons.WeapsOn.String == "3" {
		wFormOne.enableWeapon.SetChecked(true)
	} else {
		wFormOne.enableWeapon.SetChecked(false)
	}
	wFormOne.projectile.SetText(unitFunc.Missileart.String)
	if unitFunc.Missilehoming.String == "1" {
		wFormOne.projectileHoming.SetChecked(true)
	} else {
		wFormOne.projectileHoming.SetChecked(false)
	}
	wFormOne.projectileSpeed.SetText(unitFunc.Missilespeed.String)
	wFormOne.aoeFull.SetText(unit.UnitWeapons.Farea1.String)
	wFormOne.aoeMedium.SetText(unit.UnitWeapons.Harea1.String)
	wFormOne.aoeSmall.SetText(unit.UnitWeapons.Qarea1.String)
	wFormOne.weaponRange.SetText(unit.UnitWeapons.RangeN1.String)
	wFormOne.cooldown.SetText(unit.UnitWeapons.Cool1.String)
	wFormOne.damageBase.SetText(unit.UnitWeapons.Dmgplus1.String)
	wFormOne.damageDice.SetText(unit.UnitWeapons.Dice1.String)
	wFormOne.damageSides.SetText(unit.UnitWeapons.Sides1.String)
	wFormOne.damagePoint.SetText(unit.UnitWeapons.Dmgpt1.String)
	wFormOne.backswingPoint.SetText(unit.UnitWeapons.BackSw1.String)
	attackType := strings.Replace(unit.UnitWeapons.AtkType1.String, "\"", "", -1)
	var attackTypeSelected null.Int
	switch attackType {
	case "_":
		attackTypeSelected.SetValid(0)
	case "-":
		attackTypeSelected.SetValid(0)
	case "normal":
		attackTypeSelected.SetValid(1)
	case "pierce":
		attackTypeSelected.SetValid(2)
	case "siege":
		attackTypeSelected.SetValid(3)
	case "spells":
		attackTypeSelected.SetValid(4)
	case "chaos":
		attackTypeSelected.SetValid(5)
	case "magic":
		attackTypeSelected.SetValid(6)
	case "hero":
		attackTypeSelected.SetValid(7)
	}
	if attackTypeSelected.Valid {
		wFormOne.attackType.SetSelected(attackTypeSelected.Int)
	}
	weaponType := strings.Replace(unit.UnitWeapons.WeapTp1.String, "\"", "", -1)
	var weaponTypeSelected null.Int
	switch weaponType {
	case "_":
		weaponTypeSelected.SetValid(0)
	case "-":
		weaponTypeSelected.SetValid(0)
	case "normal":
		weaponTypeSelected.SetValid(1)
	case "instant":
		weaponTypeSelected.SetValid(2)
	case "artillery":
		weaponTypeSelected.SetValid(3)
	case "aline":
		weaponTypeSelected.SetValid(4)
	case "missile":
		weaponTypeSelected.SetValid(5)
	case "msplash":
		weaponTypeSelected.SetValid(6)
	case "mbounce":
		weaponTypeSelected.SetValid(7)
	case "mline":
		weaponTypeSelected.SetValid(8)
	}
	if weaponTypeSelected.Valid {
		wFormOne.weaponType.SetSelected(weaponTypeSelected.Int)
	}
	wFormOne.targets.air.SetChecked(false)
	wFormOne.targets.alive.SetChecked(false)
	wFormOne.targets.allies.SetChecked(false)
	wFormOne.targets.ancient.SetChecked(false)
	wFormOne.targets.bridge.SetChecked(false)
	wFormOne.targets.dead.SetChecked(false)
	wFormOne.targets.debris.SetChecked(false)
	wFormOne.targets.decoration.SetChecked(false)
	wFormOne.targets.enemies.SetChecked(false)
	wFormOne.targets.friend.SetChecked(false)
	wFormOne.targets.ground.SetChecked(false)
	wFormOne.targets.hero.SetChecked(false)
	wFormOne.targets.invulnerable.SetChecked(false)
	wFormOne.targets.item.SetChecked(false)
	wFormOne.targets.mechanical.SetChecked(false)
	wFormOne.targets.neutral.SetChecked(false)
	wFormOne.targets.nonancient.SetChecked(false)
	wFormOne.targets.nonhero.SetChecked(false)
	wFormOne.targets.nonsapper.SetChecked(false)
	wFormOne.targets.none.SetChecked(false)
	wFormOne.targets.notself.SetChecked(false)
	wFormOne.targets.organic.SetChecked(false)
	wFormOne.targets.player.SetChecked(false)
	wFormOne.targets.self.SetChecked(false)
	wFormOne.targets.structure.SetChecked(false)
	wFormOne.targets.sapper.SetChecked(false)
	wFormOne.targets.terrain.SetChecked(false)
	wFormOne.targets.tree.SetChecked(false)
	wFormOne.targets.vulnerable.SetChecked(false)
	wFormOne.targets.wall.SetChecked(false)
	wFormOne.targets.ward.SetChecked(false)
	targets := strings.Replace(unit.UnitWeapons.Targs1.String, "\"", "", -1)
	targetsSplit := strings.Split(targets, ",")
	for _, target := range targetsSplit {
		switch target {
		case "air":
			wFormOne.targets.air.SetChecked(true)
		case "alive":
			wFormOne.targets.alive.SetChecked(true)
		case "allies":
			wFormOne.targets.allies.SetChecked(true)
		case "ancient":
			wFormOne.targets.ancient.SetChecked(true)
		case "bridge":
			wFormOne.targets.bridge.SetChecked(true)
		case "dead":
			wFormOne.targets.dead.SetChecked(true)
		case "debris":
			wFormOne.targets.debris.SetChecked(true)
		case "decoration":
			wFormOne.targets.decoration.SetChecked(true)
		case "enemies":
			wFormOne.targets.enemies.SetChecked(true)
		case "friend":
			wFormOne.targets.friend.SetChecked(true)
		case "ground":
			wFormOne.targets.ground.SetChecked(true)
		case "hero":
			wFormOne.targets.hero.SetChecked(true)
		case "invulnerable":
			wFormOne.targets.invulnerable.SetChecked(true)
		case "item":
			wFormOne.targets.item.SetChecked(true)
		case "mechanical":
			wFormOne.targets.mechanical.SetChecked(true)
		case "neutral":
			wFormOne.targets.neutral.SetChecked(true)
		case "nonancient":
			wFormOne.targets.nonancient.SetChecked(true)
		case "nonhero":
			wFormOne.targets.nonhero.SetChecked(true)
		case "nonsapper":
			wFormOne.targets.nonsapper.SetChecked(true)
		case "none":
			wFormOne.targets.none.SetChecked(true)
		case "notself":
			wFormOne.targets.notself.SetChecked(true)
		case "organic":
			wFormOne.targets.organic.SetChecked(true)
		case "player":
			wFormOne.targets.player.SetChecked(true)
		case "self":
			wFormOne.targets.self.SetChecked(true)
		case "structure":
			wFormOne.targets.structure.SetChecked(true)
		case "sapper":
			wFormOne.targets.sapper.SetChecked(true)
		case "terrain":
			wFormOne.targets.terrain.SetChecked(true)
		case "tree":
			wFormOne.targets.tree.SetChecked(true)
		case "vulnerable":
			wFormOne.targets.vulnerable.SetChecked(true)
		case "wall":
			wFormOne.targets.wall.SetChecked(true)
		case "ward":
			wFormOne.targets.ward.SetChecked(true)
		}
	}
	wFormOne.aoeTargets.air.SetChecked(false)
	wFormOne.aoeTargets.alive.SetChecked(false)
	wFormOne.aoeTargets.allies.SetChecked(false)
	wFormOne.aoeTargets.ancient.SetChecked(false)
	wFormOne.aoeTargets.bridge.SetChecked(false)
	wFormOne.aoeTargets.dead.SetChecked(false)
	wFormOne.aoeTargets.debris.SetChecked(false)
	wFormOne.aoeTargets.decoration.SetChecked(false)
	wFormOne.aoeTargets.enemies.SetChecked(false)
	wFormOne.aoeTargets.friend.SetChecked(false)
	wFormOne.aoeTargets.ground.SetChecked(false)
	wFormOne.aoeTargets.hero.SetChecked(false)
	wFormOne.aoeTargets.invulnerable.SetChecked(false)
	wFormOne.aoeTargets.item.SetChecked(false)
	wFormOne.aoeTargets.mechanical.SetChecked(false)
	wFormOne.aoeTargets.neutral.SetChecked(false)
	wFormOne.aoeTargets.nonancient.SetChecked(false)
	wFormOne.aoeTargets.nonhero.SetChecked(false)
	wFormOne.aoeTargets.nonsapper.SetChecked(false)
	wFormOne.aoeTargets.none.SetChecked(false)
	wFormOne.aoeTargets.notself.SetChecked(false)
	wFormOne.aoeTargets.organic.SetChecked(false)
	wFormOne.aoeTargets.player.SetChecked(false)
	wFormOne.aoeTargets.self.SetChecked(false)
	wFormOne.aoeTargets.structure.SetChecked(false)
	wFormOne.aoeTargets.sapper.SetChecked(false)
	wFormOne.aoeTargets.terrain.SetChecked(false)
	wFormOne.aoeTargets.tree.SetChecked(false)
	wFormOne.aoeTargets.vulnerable.SetChecked(false)
	wFormOne.aoeTargets.wall.SetChecked(false)
	wFormOne.aoeTargets.ward.SetChecked(false)
	aoeTargets := strings.Replace(unit.UnitWeapons.SplashTargs1.String, "\"", "", -1)
	aoeTargetsSplit := strings.Split(aoeTargets, ",")
	for _, aoeTarget := range aoeTargetsSplit {
		switch aoeTarget {
		case "air":
			wFormOne.aoeTargets.air.SetChecked(true)
		case "alive":
			wFormOne.aoeTargets.alive.SetChecked(true)
		case "allies":
			wFormOne.aoeTargets.allies.SetChecked(true)
		case "ancient":
			wFormOne.aoeTargets.ancient.SetChecked(true)
		case "bridge":
			wFormOne.aoeTargets.bridge.SetChecked(true)
		case "dead":
			wFormOne.aoeTargets.dead.SetChecked(true)
		case "debris":
			wFormOne.aoeTargets.debris.SetChecked(true)
		case "decoration":
			wFormOne.aoeTargets.decoration.SetChecked(true)
		case "enemies":
			wFormOne.aoeTargets.enemies.SetChecked(true)
		case "friend":
			wFormOne.aoeTargets.friend.SetChecked(true)
		case "ground":
			wFormOne.aoeTargets.ground.SetChecked(true)
		case "hero":
			wFormOne.aoeTargets.hero.SetChecked(true)
		case "invulnerable":
			wFormOne.aoeTargets.invulnerable.SetChecked(true)
		case "item":
			wFormOne.aoeTargets.item.SetChecked(true)
		case "mechanical":
			wFormOne.aoeTargets.mechanical.SetChecked(true)
		case "neutral":
			wFormOne.aoeTargets.neutral.SetChecked(true)
		case "nonancient":
			wFormOne.aoeTargets.nonancient.SetChecked(true)
		case "nonhero":
			wFormOne.aoeTargets.nonhero.SetChecked(true)
		case "nonsapper":
			wFormOne.aoeTargets.nonsapper.SetChecked(true)
		case "none":
			wFormOne.aoeTargets.none.SetChecked(true)
		case "notself":
			wFormOne.aoeTargets.notself.SetChecked(true)
		case "organic":
			wFormOne.aoeTargets.organic.SetChecked(true)
		case "player":
			wFormOne.aoeTargets.player.SetChecked(true)
		case "self":
			wFormOne.aoeTargets.self.SetChecked(true)
		case "structure":
			wFormOne.aoeTargets.structure.SetChecked(true)
		case "sapper":
			wFormOne.aoeTargets.sapper.SetChecked(true)
		case "terrain":
			wFormOne.aoeTargets.terrain.SetChecked(true)
		case "tree":
			wFormOne.aoeTargets.tree.SetChecked(true)
		case "vulnerable":
			wFormOne.aoeTargets.vulnerable.SetChecked(true)
		case "wall":
			wFormOne.aoeTargets.wall.SetChecked(true)
		case "ward":
			wFormOne.aoeTargets.ward.SetChecked(true)
		}
	}

	// Set Weapon Form 2
	if unit.UnitWeapons.WeapsOn.Valid && unit.UnitWeapons.WeapsOn.String == "2" || unit.UnitWeapons.WeapsOn.String == "3" {
		wFormTwo.enableWeapon.SetChecked(true)
	} else {
		wFormTwo.enableWeapon.SetChecked(false)
	}
	wFormTwo.projectile.SetText(unitFunc.Missileart.String)
	if unitFunc.Missilehoming.String == "1" {
		wFormTwo.projectileHoming.SetChecked(true)
	} else {
		wFormTwo.projectileHoming.SetChecked(false)
	}
	wFormTwo.projectileSpeed.SetText(unitFunc.Missilespeed.String)
	wFormTwo.aoeFull.SetText(unit.UnitWeapons.Farea2.String)
	wFormTwo.aoeMedium.SetText(unit.UnitWeapons.Harea2.String)
	wFormTwo.aoeSmall.SetText(unit.UnitWeapons.Qarea2.String)
	wFormTwo.weaponRange.SetText(unit.UnitWeapons.RangeN2.String)
	wFormTwo.cooldown.SetText(unit.UnitWeapons.Cool2.String)
	wFormTwo.damageBase.SetText(unit.UnitWeapons.Dmgplus2.String)
	wFormTwo.damageDice.SetText(unit.UnitWeapons.Dice2.String)
	wFormTwo.damageSides.SetText(unit.UnitWeapons.Sides2.String)
	wFormTwo.damagePoint.SetText(unit.UnitWeapons.Dmgpt2.String)
	wFormTwo.backswingPoint.SetText(unit.UnitWeapons.BackSw2.String)
	attackType2 := strings.Replace(unit.UnitWeapons.AtkType2.String, "\"", "", -1)
	var attackTypeSelected2 null.Int
	switch attackType2 {
	case "_":
		attackTypeSelected2.SetValid(0)
	case "-":
		attackTypeSelected2.SetValid(0)
	case "normal":
		attackTypeSelected2.SetValid(1)
	case "pierce":
		attackTypeSelected2.SetValid(2)
	case "siege":
		attackTypeSelected2.SetValid(3)
	case "spells":
		attackTypeSelected2.SetValid(4)
	case "chaos":
		attackTypeSelected2.SetValid(5)
	case "magic":
		attackTypeSelected2.SetValid(6)
	case "hero":
		attackTypeSelected2.SetValid(7)
	}
	if attackTypeSelected2.Valid {
		wFormTwo.attackType.SetSelected(attackTypeSelected2.Int)
	}
	weaponType2 := strings.Replace(unit.UnitWeapons.WeapTp2.String, "\"", "", -1)
	var weaponTypeSelected2 null.Int
	switch weaponType2 {
	case "_":
		weaponTypeSelected2.SetValid(0)
	case "-":
		weaponTypeSelected2.SetValid(0)
	case "normal":
		weaponTypeSelected2.SetValid(1)
	case "instant":
		weaponTypeSelected2.SetValid(2)
	case "artillery":
		weaponTypeSelected2.SetValid(3)
	case "aline":
		weaponTypeSelected2.SetValid(4)
	case "missile":
		weaponTypeSelected2.SetValid(5)
	case "msplash":
		weaponTypeSelected2.SetValid(6)
	case "mbounce":
		weaponTypeSelected2.SetValid(7)
	case "mline":
		weaponTypeSelected2.SetValid(8)
	}
	if weaponTypeSelected2.Valid {
		wFormTwo.weaponType.SetSelected(weaponTypeSelected2.Int)
	}
	wFormTwo.targets.air.SetChecked(false)
	wFormTwo.targets.alive.SetChecked(false)
	wFormTwo.targets.allies.SetChecked(false)
	wFormTwo.targets.ancient.SetChecked(false)
	wFormTwo.targets.bridge.SetChecked(false)
	wFormTwo.targets.dead.SetChecked(false)
	wFormTwo.targets.debris.SetChecked(false)
	wFormTwo.targets.decoration.SetChecked(false)
	wFormTwo.targets.enemies.SetChecked(false)
	wFormTwo.targets.friend.SetChecked(false)
	wFormTwo.targets.ground.SetChecked(false)
	wFormTwo.targets.hero.SetChecked(false)
	wFormTwo.targets.invulnerable.SetChecked(false)
	wFormTwo.targets.item.SetChecked(false)
	wFormTwo.targets.mechanical.SetChecked(false)
	wFormTwo.targets.neutral.SetChecked(false)
	wFormTwo.targets.nonancient.SetChecked(false)
	wFormTwo.targets.nonhero.SetChecked(false)
	wFormTwo.targets.nonsapper.SetChecked(false)
	wFormTwo.targets.none.SetChecked(false)
	wFormTwo.targets.notself.SetChecked(false)
	wFormTwo.targets.organic.SetChecked(false)
	wFormTwo.targets.player.SetChecked(false)
	wFormTwo.targets.self.SetChecked(false)
	wFormTwo.targets.structure.SetChecked(false)
	wFormTwo.targets.sapper.SetChecked(false)
	wFormTwo.targets.terrain.SetChecked(false)
	wFormTwo.targets.tree.SetChecked(false)
	wFormTwo.targets.vulnerable.SetChecked(false)
	wFormTwo.targets.wall.SetChecked(false)
	wFormTwo.targets.ward.SetChecked(false)
	targets2 := strings.Replace(unit.UnitWeapons.Targs2.String, "\"", "", -1)
	targetsSplit2 := strings.Split(targets2, ",")
	for _, target2 := range targetsSplit2 {
		switch target2 {
		case "air":
			wFormTwo.targets.air.SetChecked(true)
		case "alive":
			wFormTwo.targets.alive.SetChecked(true)
		case "allies":
			wFormTwo.targets.allies.SetChecked(true)
		case "ancient":
			wFormTwo.targets.ancient.SetChecked(true)
		case "bridge":
			wFormTwo.targets.bridge.SetChecked(true)
		case "dead":
			wFormTwo.targets.dead.SetChecked(true)
		case "debris":
			wFormTwo.targets.debris.SetChecked(true)
		case "decoration":
			wFormTwo.targets.decoration.SetChecked(true)
		case "enemies":
			wFormTwo.targets.enemies.SetChecked(true)
		case "friend":
			wFormTwo.targets.friend.SetChecked(true)
		case "ground":
			wFormTwo.targets.ground.SetChecked(true)
		case "hero":
			wFormTwo.targets.hero.SetChecked(true)
		case "invulnerable":
			wFormTwo.targets.invulnerable.SetChecked(true)
		case "item":
			wFormTwo.targets.item.SetChecked(true)
		case "mechanical":
			wFormTwo.targets.mechanical.SetChecked(true)
		case "neutral":
			wFormTwo.targets.neutral.SetChecked(true)
		case "nonancient":
			wFormTwo.targets.nonancient.SetChecked(true)
		case "nonhero":
			wFormTwo.targets.nonhero.SetChecked(true)
		case "nonsapper":
			wFormTwo.targets.nonsapper.SetChecked(true)
		case "none":
			wFormTwo.targets.none.SetChecked(true)
		case "notself":
			wFormTwo.targets.notself.SetChecked(true)
		case "organic":
			wFormTwo.targets.organic.SetChecked(true)
		case "player":
			wFormTwo.targets.player.SetChecked(true)
		case "self":
			wFormTwo.targets.self.SetChecked(true)
		case "structure":
			wFormTwo.targets.structure.SetChecked(true)
		case "sapper":
			wFormTwo.targets.sapper.SetChecked(true)
		case "terrain":
			wFormTwo.targets.terrain.SetChecked(true)
		case "tree":
			wFormTwo.targets.tree.SetChecked(true)
		case "vulnerable":
			wFormTwo.targets.vulnerable.SetChecked(true)
		case "wall":
			wFormTwo.targets.wall.SetChecked(true)
		case "ward":
			wFormTwo.targets.ward.SetChecked(true)
		}
	}
	wFormTwo.aoeTargets.air.SetChecked(false)
	wFormTwo.aoeTargets.alive.SetChecked(false)
	wFormTwo.aoeTargets.allies.SetChecked(false)
	wFormTwo.aoeTargets.ancient.SetChecked(false)
	wFormTwo.aoeTargets.bridge.SetChecked(false)
	wFormTwo.aoeTargets.dead.SetChecked(false)
	wFormTwo.aoeTargets.debris.SetChecked(false)
	wFormTwo.aoeTargets.decoration.SetChecked(false)
	wFormTwo.aoeTargets.enemies.SetChecked(false)
	wFormTwo.aoeTargets.friend.SetChecked(false)
	wFormTwo.aoeTargets.ground.SetChecked(false)
	wFormTwo.aoeTargets.hero.SetChecked(false)
	wFormTwo.aoeTargets.invulnerable.SetChecked(false)
	wFormTwo.aoeTargets.item.SetChecked(false)
	wFormTwo.aoeTargets.mechanical.SetChecked(false)
	wFormTwo.aoeTargets.neutral.SetChecked(false)
	wFormTwo.aoeTargets.nonancient.SetChecked(false)
	wFormTwo.aoeTargets.nonhero.SetChecked(false)
	wFormTwo.aoeTargets.nonsapper.SetChecked(false)
	wFormTwo.aoeTargets.none.SetChecked(false)
	wFormTwo.aoeTargets.notself.SetChecked(false)
	wFormTwo.aoeTargets.organic.SetChecked(false)
	wFormTwo.aoeTargets.player.SetChecked(false)
	wFormTwo.aoeTargets.self.SetChecked(false)
	wFormTwo.aoeTargets.structure.SetChecked(false)
	wFormTwo.aoeTargets.sapper.SetChecked(false)
	wFormTwo.aoeTargets.terrain.SetChecked(false)
	wFormTwo.aoeTargets.tree.SetChecked(false)
	wFormTwo.aoeTargets.vulnerable.SetChecked(false)
	wFormTwo.aoeTargets.wall.SetChecked(false)
	wFormTwo.aoeTargets.ward.SetChecked(false)
	aoeTargets2 := strings.Replace(unit.UnitWeapons.SplashTargs2.String, "\"", "", -1)
	aoeTargetsSplit2 := strings.Split(aoeTargets2, ",")
	for _, aoeTarget2 := range aoeTargetsSplit2 {
		switch aoeTarget2 {
		case "air":
			wFormTwo.aoeTargets.air.SetChecked(true)
		case "alive":
			wFormTwo.aoeTargets.alive.SetChecked(true)
		case "allies":
			wFormTwo.aoeTargets.allies.SetChecked(true)
		case "ancient":
			wFormTwo.aoeTargets.ancient.SetChecked(true)
		case "bridge":
			wFormTwo.aoeTargets.bridge.SetChecked(true)
		case "dead":
			wFormTwo.aoeTargets.dead.SetChecked(true)
		case "debris":
			wFormTwo.aoeTargets.debris.SetChecked(true)
		case "decoration":
			wFormTwo.aoeTargets.decoration.SetChecked(true)
		case "enemies":
			wFormTwo.aoeTargets.enemies.SetChecked(true)
		case "friend":
			wFormTwo.aoeTargets.friend.SetChecked(true)
		case "ground":
			wFormTwo.aoeTargets.ground.SetChecked(true)
		case "hero":
			wFormTwo.aoeTargets.hero.SetChecked(true)
		case "invulnerable":
			wFormTwo.aoeTargets.invulnerable.SetChecked(true)
		case "item":
			wFormTwo.aoeTargets.item.SetChecked(true)
		case "mechanical":
			wFormTwo.aoeTargets.mechanical.SetChecked(true)
		case "neutral":
			wFormTwo.aoeTargets.neutral.SetChecked(true)
		case "nonancient":
			wFormTwo.aoeTargets.nonancient.SetChecked(true)
		case "nonhero":
			wFormTwo.aoeTargets.nonhero.SetChecked(true)
		case "nonsapper":
			wFormTwo.aoeTargets.nonsapper.SetChecked(true)
		case "none":
			wFormTwo.aoeTargets.none.SetChecked(true)
		case "notself":
			wFormTwo.aoeTargets.notself.SetChecked(true)
		case "organic":
			wFormTwo.aoeTargets.organic.SetChecked(true)
		case "player":
			wFormTwo.aoeTargets.player.SetChecked(true)
		case "self":
			wFormTwo.aoeTargets.self.SetChecked(true)
		case "structure":
			wFormTwo.aoeTargets.structure.SetChecked(true)
		case "sapper":
			wFormTwo.aoeTargets.sapper.SetChecked(true)
		case "terrain":
			wFormTwo.aoeTargets.terrain.SetChecked(true)
		case "tree":
			wFormTwo.aoeTargets.tree.SetChecked(true)
		case "vulnerable":
			wFormTwo.aoeTargets.vulnerable.SetChecked(true)
		case "wall":
			wFormTwo.aoeTargets.wall.SetChecked(true)
		case "ward":
			wFormTwo.aoeTargets.ward.SetChecked(true)
		}
	}

	// Set Data Form
	builds := strings.Replace(unitFunc.Builds.String, "\"", "", -1)
	dForm.acquisition.SetText(unit.UnitWeapons.Acquire.String)
	dForm.builds.SetText(builds)
	dForm.health.SetText(unit.UnitBalance.HP.String)
	dForm.mana.SetText(unit.UnitBalance.ManaN.String)
	if unit.UnitBalance.Isbldg.String == "1" {
		dForm.isBuilding.SetChecked(true)
	} else {
		dForm.isBuilding.SetChecked(false)
	}
	dForm.defense.SetText(unit.UnitBalance.Def.String)
	dForm.lumberCost.SetText(unit.UnitBalance.Lumbercost.String)
	dForm.goldCost.SetText(unit.UnitBalance.Goldcost.String)
	dForm.foodCost.SetText(unit.UnitBalance.Fused.String)
	dForm.foodProduction.SetText(unit.UnitBalance.Fmade.String)
	dForm.movementSpeed.SetText(unit.UnitBalance.Spd.String)
	dForm.flyingHeight.SetText(unit.UnitData.MoveHeight.String)
	dForm.minimumFlyingHeight.SetText(unit.UnitData.MoveFloor.String)
	race := strings.Replace(unit.UnitData.Race.String, "\"", "", -1)
	var raceSelected null.Int
	switch race {
	case "_":
		raceSelected.SetValid(0)
	case "-":
		raceSelected.SetValid(0)
	case "commoner":
		raceSelected.SetValid(1)
	case "creeps":
		raceSelected.SetValid(2)
	case "critters":
		raceSelected.SetValid(3)
	case "demon":
		raceSelected.SetValid(4)
	case "human":
		raceSelected.SetValid(5)
	case "naga":
		raceSelected.SetValid(6)
	case "nightelf":
		raceSelected.SetValid(7)
	case "orc":
		raceSelected.SetValid(8)
	case "other":
		raceSelected.SetValid(9)
	case "unknown":
		raceSelected.SetValid(9)
	case "undead":
		raceSelected.SetValid(10)
	}
	if raceSelected.Valid {
		dForm.race.SetSelected(raceSelected.Int)
	}
	defenseType := strings.Replace(unit.UnitBalance.DefType.String, "\"", "", -1)
	var defenseTypeSelected null.Int
	switch defenseType {
	case "normal":
		defenseTypeSelected.SetValid(0)
	case "small":
		defenseTypeSelected.SetValid(1)
	case "medium":
		defenseTypeSelected.SetValid(2)
	case "large":
		defenseTypeSelected.SetValid(3)
	case "fort":
		defenseTypeSelected.SetValid(4)
	case "hero":
		defenseTypeSelected.SetValid(5)
	case "divine":
		defenseTypeSelected.SetValid(6)
	case "none":
		defenseTypeSelected.SetValid(7)

	}
	if defenseTypeSelected.Valid {
		dForm.defenseType.SetSelected(defenseTypeSelected.Int)
	}
	movementType := strings.Replace(unit.UnitData.Movetp.String, "\"", "", -1)
	var movementTypeSelected null.Int
	switch movementType {
	case "_":
		movementTypeSelected.SetValid(0)
	case "-":
		movementTypeSelected.SetValid(0)
	case "foot":
		movementTypeSelected.SetValid(1)
	case "horse":
		movementTypeSelected.SetValid(2)
	case "fly":
		movementTypeSelected.SetValid(3)
	case "hover":
		movementTypeSelected.SetValid(4)
	case "float":
		movementTypeSelected.SetValid(5)
	case "amph":
		movementTypeSelected.SetValid(6)
	}
	if movementTypeSelected.Valid {
		dForm.movementType.SetSelected(movementTypeSelected.Int)
	}

	// Set Other Form
	if unit.UnitUI.DropItems.Valid && unit.UnitUI.DropItems.String == "1" {
		oForm.dropsItems.SetChecked(true)
	} else {
		oForm.dropsItems.SetChecked(false)
	}
	if unit.UnitUI.ElevPts.Valid {
		oForm.elevationSamplePoints.SetText(unit.UnitUI.ElevPts.String)
	}
	if unit.UnitUI.ElevRad.Valid {
		oForm.elevationSampleRadius.SetText(unit.UnitUI.ElevRad.String)
	}
	if unit.UnitData.IsBuildOn.Valid && unit.UnitData.IsBuildOn.String == "1" {
		oForm.canBeBuiltOn.SetChecked(true)
	} else {
		oForm.canBeBuiltOn.SetChecked(false)
	}
	if unit.UnitData.CanBuildOn.Valid && unit.UnitData.CanBuildOn.String == "1" {
		oForm.canBuildOn.SetChecked(true)
	} else {
		oForm.canBuildOn.SetChecked(false)
	}
	if unit.UnitData.TurnRate.Valid {
		oForm.turnRate.SetText(unit.UnitData.TurnRate.String)
	}
	if unit.UnitData.CargoSize.Valid {
		oForm.cargoSize.SetText(unit.UnitData.CargoSize.String)
	}
	if unit.UnitData.CanSleep.Valid && unit.UnitData.CanSleep.String == "1" {
		oForm.canSleep.SetChecked(true)
	} else {
		oForm.canSleep.SetChecked(false)
	}
	if unit.UnitData.Death.Valid {
		oForm.death.SetText(unit.UnitData.Death.String)
	}
	oForm.targetedAs.air.SetChecked(false)
	oForm.targetedAs.alive.SetChecked(false)
	oForm.targetedAs.allies.SetChecked(false)
	oForm.targetedAs.ancient.SetChecked(false)
	oForm.targetedAs.bridge.SetChecked(false)
	oForm.targetedAs.dead.SetChecked(false)
	oForm.targetedAs.debris.SetChecked(false)
	oForm.targetedAs.decoration.SetChecked(false)
	oForm.targetedAs.enemies.SetChecked(false)
	oForm.targetedAs.friend.SetChecked(false)
	oForm.targetedAs.ground.SetChecked(false)
	oForm.targetedAs.hero.SetChecked(false)
	oForm.targetedAs.invulnerable.SetChecked(false)
	oForm.targetedAs.item.SetChecked(false)
	oForm.targetedAs.mechanical.SetChecked(false)
	oForm.targetedAs.neutral.SetChecked(false)
	oForm.targetedAs.nonancient.SetChecked(false)
	oForm.targetedAs.nonhero.SetChecked(false)
	oForm.targetedAs.nonsapper.SetChecked(false)
	oForm.targetedAs.none.SetChecked(false)
	oForm.targetedAs.notself.SetChecked(false)
	oForm.targetedAs.organic.SetChecked(false)
	oForm.targetedAs.player.SetChecked(false)
	oForm.targetedAs.self.SetChecked(false)
	oForm.targetedAs.structure.SetChecked(false)
	oForm.targetedAs.sapper.SetChecked(false)
	oForm.targetedAs.terrain.SetChecked(false)
	oForm.targetedAs.tree.SetChecked(false)
	oForm.targetedAs.vulnerable.SetChecked(false)
	oForm.targetedAs.wall.SetChecked(false)
	oForm.targetedAs.ward.SetChecked(false)
	targetedAs := strings.Replace(unit.UnitData.TargType.String, "\"", "", -1)
	targetedAsSplit := strings.Split(targetedAs, ",")
	for _, targeted := range targetedAsSplit {
		switch targeted {
		case "air":
			oForm.targetedAs.air.SetChecked(true)
		case "alive":
			oForm.targetedAs.alive.SetChecked(true)
		case "allies":
			oForm.targetedAs.allies.SetChecked(true)
		case "ancient":
			oForm.targetedAs.ancient.SetChecked(true)
		case "bridge":
			oForm.targetedAs.bridge.SetChecked(true)
		case "dead":
			oForm.targetedAs.dead.SetChecked(true)
		case "debris":
			oForm.targetedAs.debris.SetChecked(true)
		case "decoration":
			oForm.targetedAs.decoration.SetChecked(true)
		case "enemies":
			oForm.targetedAs.enemies.SetChecked(true)
		case "friend":
			oForm.targetedAs.friend.SetChecked(true)
		case "ground":
			oForm.targetedAs.ground.SetChecked(true)
		case "hero":
			oForm.targetedAs.hero.SetChecked(true)
		case "invulnerable":
			oForm.targetedAs.invulnerable.SetChecked(true)
		case "item":
			oForm.targetedAs.item.SetChecked(true)
		case "mechanical":
			oForm.targetedAs.mechanical.SetChecked(true)
		case "neutral":
			oForm.targetedAs.neutral.SetChecked(true)
		case "nonancient":
			oForm.targetedAs.nonancient.SetChecked(true)
		case "nonhero":
			oForm.targetedAs.nonhero.SetChecked(true)
		case "nonsapper":
			oForm.targetedAs.nonsapper.SetChecked(true)
		case "none":
			oForm.targetedAs.none.SetChecked(true)
		case "notself":
			oForm.targetedAs.notself.SetChecked(true)
		case "organic":
			oForm.targetedAs.organic.SetChecked(true)
		case "player":
			oForm.targetedAs.player.SetChecked(true)
		case "self":
			oForm.targetedAs.self.SetChecked(true)
		case "structure":
			oForm.targetedAs.structure.SetChecked(true)
		case "sapper":
			oForm.targetedAs.sapper.SetChecked(true)
		case "terrain":
			oForm.targetedAs.terrain.SetChecked(true)
		case "tree":
			oForm.targetedAs.tree.SetChecked(true)
		case "vulnerable":
			oForm.targetedAs.vulnerable.SetChecked(true)
		case "wall":
			oForm.targetedAs.wall.SetChecked(true)
		case "ward":
			oForm.targetedAs.ward.SetChecked(true)
		}
	}
}

func makePathingTextureComboBox() *ui.Combobox {
	comboBox := ui.NewCombobox()
	comboBox.Append("NONE")
	comboBox.Append("Altar of Darkness")
	comboBox.Append("Altar of Elders")
	comboBox.Append("Altar of Kings")
	comboBox.Append("Altar of Storms")
	comboBox.Append("Ancient of Lore")
	comboBox.Append("Ancient Protector")
	comboBox.Append("Arcane Tower")
	comboBox.Append("Castle")
	comboBox.Append("Demon Gate")
	comboBox.Append("Dimensional Gate (facing southeast)")
	comboBox.Append("Dimensional Gate (facing southwest)")
	comboBox.Append("Gold Mine")
	comboBox.Append("Fortress")
	comboBox.Append("Halls of the Dead")
	comboBox.Append("Keep")
	comboBox.Append("Moon Well")
	comboBox.Append("Orc Shipyard")
	comboBox.Append("Spawning Grounds")
	comboBox.Append("Spirit Tower")
	comboBox.Append("Temple of Tides")

	comboBox.SetSelected(0)

	return comboBox
}

func saveUnitsToFile() {
	customUnitFuncs := new(models.UnitFuncs)
	campaignUnitFuncs := make([]*models.UnitFunc, len(unitFuncMap))
	var campaignIndex = 0
	for _, k := range unitFuncMap {
		campaignUnitFuncs[campaignIndex] = k
		campaignIndex++
	}

	unitMapLength := len(baseUnitMap)
	parsedSLKUnitsAbilities := make([]*models.UnitAbilities, unitMapLength)
	parsedSLKUnitsData := make([]*models.UnitData, unitMapLength)
	parsedSLKUnitsUI := make([]*models.UnitUI, unitMapLength)
	parsedSLKUnitsWeapons := make([]*models.UnitWeapons, unitMapLength)
	parsedSLKUnitsBalance := make([]*models.UnitBalance, unitMapLength)

	var i = 0
	for _, parsedSLKUnit := range baseUnitMap {
		parsedSLKUnitsAbilities[i] = parsedSLKUnit.UnitAbilities
		parsedSLKUnitsData[i] = parsedSLKUnit.UnitData
		parsedSLKUnitsUI[i] = parsedSLKUnit.UnitUI
		parsedSLKUnitsWeapons[i] = parsedSLKUnit.UnitWeapons
		parsedSLKUnitsBalance[i] = parsedSLKUnit.UnitBalance
		i++
	}

	parser.WriteToFiles(customUnitFuncs, parsedSLKUnitsAbilities, parsedSLKUnitsData, parsedSLKUnitsUI, parsedSLKUnitsWeapons, parsedSLKUnitsBalance)
}

func makeRaceComboBox() *ui.Combobox {
	comboBox := ui.NewCombobox()
	comboBox.Append("None")
	comboBox.Append("Commoner")
	comboBox.Append("Creep")
	comboBox.Append("Critter")
	comboBox.Append("Demon")
	comboBox.Append("Human")
	comboBox.Append("Naga")
	comboBox.Append("Night Elf")
	comboBox.Append("Orc")
	comboBox.Append("Other")
	comboBox.Append("Undead")

	comboBox.SetSelected(0)

	return comboBox
}

func makeMovementTypeComboBox() *ui.Combobox {
	comboBox := ui.NewCombobox()
	comboBox.Append("NONE")
	comboBox.Append("Foot")
	comboBox.Append("Horse")
	comboBox.Append("Fly")
	comboBox.Append("Hover")
	comboBox.Append("Float")
	comboBox.Append("Amphibious")

	comboBox.SetSelected(0)

	return comboBox
}

func makeDefenseTypeComboBox() *ui.Combobox {
	comboBox := ui.NewCombobox()
	comboBox.Append("Normal")
	comboBox.Append("Small")
	comboBox.Append("Medium")
	comboBox.Append("Large")
	comboBox.Append("Fortified")
	comboBox.Append("Hero")
	comboBox.Append("Divine")
	comboBox.Append("Unarmored")

	comboBox.SetSelected(0)

	return comboBox
}

func makeAttackTypeComboBox() *ui.Combobox {
	comboBox := ui.NewCombobox()
	comboBox.Append("None")
	comboBox.Append("Normal")
	comboBox.Append("Pierce")
	comboBox.Append("Siege")
	comboBox.Append("Spells")
	comboBox.Append("Chaos")
	comboBox.Append("Magic")
	comboBox.Append("Hero")

	comboBox.SetSelected(0)

	return comboBox
}

func makeWeaponTypeComboBox() *ui.Combobox {
	comboBox := ui.NewCombobox()
	comboBox.Append("NONE")
	comboBox.Append("Normal")
	comboBox.Append("Instant")
	comboBox.Append("Artillery")
	comboBox.Append("Artillery (Line)")
	comboBox.Append("Missile")
	comboBox.Append("Missile (Splash)")
	comboBox.Append("Missile (Bounce)")
	comboBox.Append("Missile (Line)")

	comboBox.SetSelected(0)

	return comboBox
}

func makeTargetTypeGrid() targetGrid {
	tGrid := targetGrid{ui.NewGrid(), ui.NewCheckbox("air"), ui.NewCheckbox("alive"), ui.NewCheckbox("allies"), ui.NewCheckbox("ancient"), ui.NewCheckbox("bridge"), ui.NewCheckbox("dead"), ui.NewCheckbox("debris"), ui.NewCheckbox("decoration"), ui.NewCheckbox("enemies"), ui.NewCheckbox("friend"), ui.NewCheckbox("ground"), ui.NewCheckbox("hero"), ui.NewCheckbox("invulnerable"), ui.NewCheckbox("item"), ui.NewCheckbox("mechanical"), ui.NewCheckbox("neutral"), ui.NewCheckbox("nonancient"), ui.NewCheckbox("nonhero"), ui.NewCheckbox("nonsapper"), ui.NewCheckbox("none"), ui.NewCheckbox("notself"), ui.NewCheckbox("organic"), ui.NewCheckbox("player"), ui.NewCheckbox("self"), ui.NewCheckbox("structure"), ui.NewCheckbox("sapper"), ui.NewCheckbox("terrain"), ui.NewCheckbox("tree"), ui.NewCheckbox("vulnerable"), ui.NewCheckbox("wall"), ui.NewCheckbox("ward")}

	tGrid.Append(tGrid.air, 0, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.alive, 1, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.allies, 2, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.ancient, 3, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.bridge, 4, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.dead, 5, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.debris, 0, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.decoration, 1, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.enemies, 2, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.friend, 3, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.ground, 4, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.hero, 5, 1, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.invulnerable, 0, 2, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.item, 1, 2, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.mechanical, 2, 2, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.neutral, 3, 2, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.nonancient, 4, 2, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.nonhero, 5, 2, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.nonsapper, 0, 3, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.none, 1, 3, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.notself, 2, 3, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.organic, 3, 3, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.player, 4, 3, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.self, 5, 3, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.structure, 0, 4, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.sapper, 1, 4, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.terrain, 2, 4, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.tree, 3, 4, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.vulnerable, 4, 4, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.wall, 5, 4, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	tGrid.Append(tGrid.ward, 0, 5, 1, 1, false, ui.AlignFill, false, ui.AlignFill)

	return tGrid
}

func makeUnitInputForm() *ui.Tab {
	tab := ui.NewTab()

	uForm = uiForm{ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewCheckbox(""), ui.NewEntry(), ui.NewEntry(), makePathingTextureComboBox(), ui.NewColorButton(), ui.NewSlider(0, 255), ui.NewSlider(0, 255), ui.NewSlider(0, 255), ui.NewEntry(), ui.NewEntry(), ui.NewMultilineEntry()}
	uiForm := ui.NewForm()
	uiForm.SetPadded(true)

	uForm.color.OnChanged(func(button *ui.ColorButton) {
		r, g, b, _ := button.Color()
		uForm.red.SetValue(int(255 * r))
		uForm.green.SetValue(int(255 * g))
		uForm.blue.SetValue(int(255 * b))
	})
	uForm.red.OnChanged(func(slider *ui.Slider) {
		uForm.color.SetColor(float64(uForm.red.Value())*0.003921569, float64(uForm.green.Value())*0.003921569, float64(uForm.blue.Value())*0.003921569, 1)
	})
	uForm.green.OnChanged(func(slider *ui.Slider) {
		uForm.color.SetColor(float64(uForm.red.Value())*0.003921569, float64(uForm.green.Value())*0.003921569, float64(uForm.blue.Value())*0.003921569, 1)
	})
	uForm.blue.OnChanged(func(slider *ui.Slider) {
		uForm.color.SetColor(float64(uForm.red.Value())*0.003921569, float64(uForm.green.Value())*0.003921569, float64(uForm.blue.Value())*0.003921569, 1)
	})

	uForm.red.SetValue(255)
	uForm.green.SetValue(255)
	uForm.blue.SetValue(255)
	uForm.color.SetColor(1, 1, 1, 1)

	generateUnitIdButton := ui.NewButton("Generate Valid ID")
	generateUnitIdButton.OnClicked(func(button *ui.Button) {
		uForm.unitId.SetText(getNextValidUnitId(lastValidIndex))
	})
	hBox := ui.NewHorizontalBox()
	hBox.Append(uForm.unitId, true)
	hBox.Append(generateUnitIdButton, false)

	uiForm.Append("UnitID", hBox, false)
	uiForm.Append("Name", uForm.name, false)
	uiForm.Append("Abilities", uForm.abilities, false)
	uiForm.Append("Icon", uForm.icon, false)
	uiForm.Append("Button X", uForm.buttonX, false)
	uiForm.Append("Button Y", uForm.buttonY, false)
	uiForm.Append("Model", uForm.model, false)
	uiForm.Append("Sound Set", uForm.soundSet, false)
	uiForm.Append("Hide Minimap Display", uForm.hideMinimapDisplay, false)
	uiForm.Append("Scaling Value", uForm.scalingValue, false)
	uiForm.Append("Selection Scale", uForm.selectionScale, false)
	uiForm.Append("Pathing Texture", uForm.pathingTexture, false)
	uiForm.Append("Color", uForm.color, false)
	uiForm.Append("Red", uForm.red, false)
	uiForm.Append("Green", uForm.green, false)
	uiForm.Append("Blue", uForm.blue, false)
	uiForm.Append("Hotkey", uForm.hotkey, false)
	uiForm.Append("Tooltip", uForm.tooltip, false)
	uiForm.Append("Description", uForm.description, true)

	tab.Append("UI", uiForm)
	tab.SetMargined(0, true)

	wFormOne = weaponForm{ui.NewCheckbox(""), ui.NewEntry(), ui.NewEntry(), makeAttackTypeComboBox(), makeTargetTypeGrid(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), makeWeaponTypeComboBox(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), makeTargetTypeGrid(), ui.NewEntry(), ui.NewCheckbox(""), ui.NewEntry()}
	wFormTwo = weaponForm{ui.NewCheckbox(""), ui.NewEntry(), ui.NewEntry(), makeAttackTypeComboBox(), makeTargetTypeGrid(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), makeWeaponTypeComboBox(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), makeTargetTypeGrid(), ui.NewEntry(), ui.NewCheckbox(""), ui.NewEntry()}

	weaponTab := weaponTab{ui.NewTab(), wFormOne, wFormTwo}

	weaponFormOne := ui.NewForm()
	weaponFormOne.SetPadded(true)

	areaOfEffectHBoxOne := ui.NewHorizontalBox()
	aoeFullFormOne := ui.NewForm()
	aoeFullFormOne.SetPadded(true)
	aoeFullFormOne.Append("Full", wFormOne.aoeFull, false)
	aoeMediumFormOne := ui.NewForm()
	aoeMediumFormOne.SetPadded(true)
	aoeMediumFormOne.Append("Medium", wFormOne.aoeMedium, false)
	aoeSmallFormOne := ui.NewForm()
	aoeSmallFormOne.SetPadded(true)
	aoeSmallFormOne.Append("Small", wFormOne.aoeSmall, false)
	areaOfEffectHBoxOne.Append(aoeFullFormOne, false)
	areaOfEffectHBoxOne.Append(aoeMediumFormOne, false)
	areaOfEffectHBoxOne.Append(aoeSmallFormOne, false)

	weaponFormOne.Append("Enable Weapon", wFormOne.enableWeapon, false)
	weaponFormOne.Append("Backswing Point", wFormOne.backswingPoint, false)
	weaponFormOne.Append("Damage Point", wFormOne.damagePoint, false)
	weaponFormOne.Append("Attack Type", wFormOne.attackType, false)
	weaponFormOne.Append("Targets", wFormOne.targets, false)
	weaponFormOne.Append("Cooldown", wFormOne.cooldown, false)
	weaponFormOne.Append("Damage Base", wFormOne.damageBase, false)
	weaponFormOne.Append("Damage Dice", wFormOne.damageDice, false)
	weaponFormOne.Append("Damage Sides", wFormOne.damageSides, false)
	weaponFormOne.Append("Range", wFormOne.weaponRange, false)
	weaponFormOne.Append("Weapon Type", wFormOne.weaponType, false)
	weaponFormOne.Append("AOE", areaOfEffectHBoxOne, false)
	weaponFormOne.Append("AOE Targets", wFormOne.aoeTargets, false)
	weaponFormOne.Append("Projectile", wFormOne.projectile, false)
	weaponFormOne.Append("Projectile Homing", wFormOne.projectileHoming, false)
	weaponFormOne.Append("Projectile Speed", wFormOne.projectileSpeed, false)

	weaponTab.Append("Weapon 1", weaponFormOne)
	weaponTab.SetMargined(0, true)

	weaponFormTwo := ui.NewForm()
	weaponFormTwo.SetPadded(true)

	areaOfEffectHBoxTwo := ui.NewHorizontalBox()
	aoeFullFormTwo := ui.NewForm()
	aoeFullFormTwo.SetPadded(true)
	aoeFullFormTwo.Append("Full", wFormTwo.aoeFull, false)
	aoeMediumFormTwo := ui.NewForm()
	aoeMediumFormTwo.SetPadded(true)
	aoeMediumFormTwo.Append("Medium", wFormTwo.aoeMedium, false)
	aoeSmallFormTwo := ui.NewForm()
	aoeSmallFormTwo.SetPadded(true)
	aoeSmallFormTwo.Append("Small", wFormTwo.aoeSmall, false)
	areaOfEffectHBoxTwo.Append(aoeFullFormTwo, false)
	areaOfEffectHBoxTwo.Append(aoeMediumFormTwo, false)
	areaOfEffectHBoxTwo.Append(aoeSmallFormTwo, false)

	weaponFormTwo.Append("Enable Weapon", wFormTwo.enableWeapon, false)
	weaponFormTwo.Append("Backswing Point", wFormTwo.backswingPoint, false)
	weaponFormTwo.Append("Damage Point", wFormTwo.damagePoint, false)
	weaponFormTwo.Append("Attack Type", wFormTwo.attackType, false)
	weaponFormTwo.Append("Targets", wFormTwo.targets, false)
	weaponFormTwo.Append("Cooldown", wFormTwo.cooldown, false)
	weaponFormTwo.Append("Damage Base", wFormTwo.damageBase, false)
	weaponFormTwo.Append("Damage Dice", wFormTwo.damageDice, false)
	weaponFormTwo.Append("Damage Sides", wFormTwo.damageSides, false)
	weaponFormTwo.Append("Range", wFormTwo.weaponRange, false)
	weaponFormTwo.Append("Weapon Type", wFormTwo.weaponType, false)
	weaponFormTwo.Append("AOE", areaOfEffectHBoxTwo, false)
	weaponFormTwo.Append("AOE Targets", wFormTwo.aoeTargets, false)
	weaponFormTwo.Append("Projectile", wFormTwo.projectile, false)
	weaponFormTwo.Append("Projectile Homing", wFormTwo.projectileHoming, false)
	weaponFormTwo.Append("Projectile Speed", wFormTwo.projectileSpeed, false)

	weaponTab.Append("Weapon 2", weaponFormTwo)
	weaponTab.SetMargined(1, true)

	tab.Append("Weapons", weaponTab)
	tab.SetMargined(1, true)

	isBuildingCheckbox := ui.NewCheckbox("")
	isBuildingCheckbox.SetChecked(true)

	isCategorizedCampaignCheckbox := ui.NewCheckbox("")
	isCategorizedCampaignCheckbox.SetChecked(true)
	isCategorizedCampaignCheckbox.Disable()

	dForm = dataForm{ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), isBuildingCheckbox, ui.NewEntry(), makeDefenseTypeComboBox(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), makeRaceComboBox(), ui.NewEntry(), ui.NewEntry(), isCategorizedCampaignCheckbox, makeMovementTypeComboBox(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry()}

	dataForm := ui.NewForm()
	dataForm.SetPadded(true)

	movementHBox := ui.NewHorizontalBox()
	movementSpeedForm := ui.NewForm()
	movementSpeedForm.SetPadded(true)
	movementSpeedForm.Append("Speed", dForm.movementSpeed, false)
	movementSpeedMinimumForm := ui.NewForm()
	movementSpeedMinimumForm.SetPadded(true)
	movementSpeedMinimumForm.Append("Minimum", dForm.movementSpeedMinimum, false)
	movementSpeedMaximumForm := ui.NewForm()
	movementSpeedMaximumForm.SetPadded(true)
	movementSpeedMaximumForm.Append("Maximum", dForm.movementSpeedMaximum, false)
	movementHBox.Append(movementSpeedForm, false)
	movementHBox.Append(movementSpeedMinimumForm, false)
	movementHBox.Append(movementSpeedMaximumForm, false)

	dataForm.Append("Acquisition Range", dForm.acquisition, false)
	dataForm.Append("Builds", dForm.builds, false)
	dataForm.Append("Sells", dForm.sells, false)
	dataForm.Append("Upgrades To", dForm.upgradesTo, false)
	dataForm.Append("Trains", dForm.trains, false)
	dataForm.Append("Health", dForm.health, false)
	dataForm.Append("Mana", dForm.mana, false)
	dataForm.Append("Is Building", dForm.isBuilding, false)
	dataForm.Append("Defense", dForm.defense, false)
	dataForm.Append("Defense Type", dForm.defenseType, false)
	dataForm.Append("Lumber Cost", dForm.lumberCost, false)
	dataForm.Append("Gold Cost", dForm.goldCost, false)
	dataForm.Append("Points", dForm.points, false)
	dataForm.Append("Race", dForm.race, false)
	dataForm.Append("Food Cost", dForm.foodCost, false)
	dataForm.Append("Food Production", dForm.foodProduction, false)
	dataForm.Append("Is Campaign", dForm.isCampaign, false)
	dataForm.Append("Movement Type", dForm.movementType, false)
	dataForm.Append("Movement", movementHBox, false)
	dataForm.Append("Flying Height", dForm.flyingHeight, false)
	dataForm.Append("Flying Height (Minimum)", dForm.minimumFlyingHeight, false)

	tab.Append("Data", dataForm)
	tab.SetMargined(2, true)

	oForm = otherForm{ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewEntry(), ui.NewCheckbox(""), ui.NewCheckbox(""), ui.NewCheckbox(""), ui.NewCheckbox(""), ui.NewEntry(), ui.NewEntry(), makeTargetTypeGrid(), ui.NewEntry()}

	otherForm := ui.NewForm()
	otherForm.SetPadded(true)

	// otherForm.Append("Death Type", oForm.deathType, false)
	otherForm.Append("Death", oForm.death, false)
	otherForm.Append("Cargo Size", oForm.cargoSize, false)
	otherForm.Append("Turn Rate", oForm.turnRate, false)
	otherForm.Append("Can Sleep", oForm.canSleep, false)
	otherForm.Append("Can Be Built On", oForm.canBeBuiltOn, false)
	otherForm.Append("Can Build On", oForm.canBuildOn, false)
	otherForm.Append("Drops Items", oForm.dropsItems, false)
	otherForm.Append("Elevation Sample Points", oForm.elevationSamplePoints, false)
	otherForm.Append("Elevation Sample Radius", oForm.elevationSampleRadius, false)
	otherForm.Append("Targeted As", oForm.targetedAs, false)
	otherForm.Append("Level", oForm.level, false)

	tab.Append("Other", otherForm)
	tab.SetMargined(3, true)

	return tab
}

func targetsAppendString(baseString string, appendString string) string {
	if baseString != "" {
		return baseString + "," + appendString
	} else {
		return appendString
	}
}

func targetsToString(grid targetGrid) string {
	str := ""
	if grid.air.Checked() {
		str = targetsAppendString(str, "air")
	}
	if grid.alive.Checked() {
		str = targetsAppendString(str, "alive")
	}
	if grid.allies.Checked() {
		str = targetsAppendString(str, "allies")
	}
	if grid.ancient.Checked() {
		str = targetsAppendString(str, "ancient")
	}
	if grid.bridge.Checked() {
		str = targetsAppendString(str, "bridge")
	}
	if grid.dead.Checked() {
		str = targetsAppendString(str, "dead")
	}
	if grid.debris.Checked() {
		str = targetsAppendString(str, "debris")
	}
	if grid.decoration.Checked() {
		str = targetsAppendString(str, "decoration")
	}
	if grid.enemies.Checked() {
		str = targetsAppendString(str, "enemies")
	}
	if grid.friend.Checked() {
		str = targetsAppendString(str, "friend")
	}
	if grid.ground.Checked() {
		str = targetsAppendString(str, "ground")
	}
	if grid.hero.Checked() {
		str = targetsAppendString(str, "hero")
	}
	if grid.invulnerable.Checked() {
		str = targetsAppendString(str, "invulnerable")
	}
	if grid.item.Checked() {
		str = targetsAppendString(str, "item")
	}
	if grid.mechanical.Checked() {
		str = targetsAppendString(str, "mechanical")
	}
	if grid.neutral.Checked() {
		str = targetsAppendString(str, "neutral")
	}
	if grid.nonancient.Checked() {
		str = targetsAppendString(str, "nonancient")
	}
	if grid.none.Checked() {
		str = targetsAppendString(str, "none")
	}
	if grid.nonhero.Checked() {
		str = targetsAppendString(str, "nonhero")
	}
	if grid.nonsapper.Checked() {
		str = targetsAppendString(str, "nonsapper")
	}
	if grid.notself.Checked() {
		str = targetsAppendString(str, "notself")
	}
	if grid.organic.Checked() {
		str = targetsAppendString(str, "organic")
	}
	if grid.player.Checked() {
		str = targetsAppendString(str, "player")
	}
	if grid.sapper.Checked() {
		str = targetsAppendString(str, "sapper")
	}
	if grid.self.Checked() {
		str = targetsAppendString(str, "self")
	}
	if grid.structure.Checked() {
		str = targetsAppendString(str, "structure")
	}
	if grid.terrain.Checked() {
		str = targetsAppendString(str, "terrain")
	}
	if grid.tree.Checked() {
		str = targetsAppendString(str, "tree")
	}
	if grid.vulnerable.Checked() {
		str = targetsAppendString(str, "vulnerable")
	}
	if grid.wall.Checked() {
		str = targetsAppendString(str, "wall")
	}
	if grid.ward.Checked() {
		str = targetsAppendString(str, "ward")
	}

	return str
}

func makeBasicControlsPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	model := ui.NewTableModel(mh)
	table := ui.NewTable(&ui.TableParams{
		Model:                         model,
		RowBackgroundColorModelColumn: 1,
	})

	table.AppendButtonColumn("Units",
		0, ui.TableModelColumnAlwaysEditable)

	tableVbox := ui.NewVerticalBox()
	/*
	saveToFileButton := ui.NewButton("Save File")
	saveToFileButton.OnClicked(func(button *ui.Button) {

	})
	tableVbox.Append(saveToFileButton, false)
	*/
	searchEntry := ui.NewSearchEntry()
	searchEntry.OnChanged(func(entry *ui.Entry) {
		var searchRegex = regexp.MustCompile(strings.ToLower(entry.Text()))
		var newList []string

		for range mh.slkUnitIdList {
			model.RowDeleted(0)
		}

		mh.rows = 0
		mh.slkUnitIdList = []string{}

		var i = 0
		for key, value := range unitFuncMap {
			if searchRegex.MatchString(strings.ToLower(key + " - " + value.Name.String)) {
				newList = append(newList, key)
				model.RowInserted(i)
				i++
			}
		}

		mh.slkUnitIdList = newList
		mh.rows = len(newList)
	})
	removeButton := ui.NewButton("Remove Unit")
	removeButton.OnClicked(func(button *ui.Button) {
		if selectedUnit.Valid {
			delete(unitFuncMap, mh.slkUnitIdList[selectedUnit.Int])
			delete(baseUnitMap, mh.slkUnitIdList[selectedUnit.Int])
			mh.slkUnitIdList = append(mh.slkUnitIdList[:selectedUnit.Int], mh.slkUnitIdList[selectedUnit.Int+1:]...)
			mh.rows = mh.rows - 1
			model.RowDeleted(selectedUnit.Int)

			selectedUnit.Int = 0
			selectedUnit.Valid = false
		}
	})
	tableVbox.Append(searchEntry, false)
	tableVbox.Append(table, true)
	tableVbox.Append(removeButton, false)
	hbox.Append(tableVbox, true)
	formVbox := ui.NewVerticalBox()
	formVbox.Append(makeUnitInputForm(), true)
	saveUnitButton := ui.NewButton("Save Unit")
	saveUnitButton.OnClicked(func(button *ui.Button) {
		if uForm.unitId.Text() == "" || len(uForm.unitId.Text()) != 4 {
			log.Println("Error: Unit ID has to be 4 characters long!")
			return
		}

		newUnitFunc := new(models.UnitFunc)
		newSlkUnit := new(models.SLKUnit)
		unitWeapons := new(models.UnitWeapons)
		unitData := new(models.UnitData)
		unitBalance := new(models.UnitBalance)
		unitUi := new(models.UnitUI)
		unitAbilities := new(models.UnitAbilities)

		newUnitFunc.UnitId = uForm.unitId.Text()
		if dForm.builds.Text() != "" {
			newUnitFunc.Builds.SetValid(dForm.builds.Text())
		}
		if uForm.name.Text() != "" {
			newUnitFunc.Name.SetValid(uForm.name.Text())
		}
		var missileSpeed = ""
		var missileHoming = ""
		var missileArt = ""
		if wFormOne.weaponType.Selected() != 0 && wFormOne.weaponType.Selected() != 1 { // 0 == None, 1 == Normal
			if wFormOne.projectileSpeed.Text() != "" {
				missileSpeed = wFormOne.projectileSpeed.Text()
			}
			if wFormOne.projectileHoming.Checked() {
				missileHoming = "1"
			} else {
				missileHoming = "0"
			}
			if wFormOne.projectile.Text() != "" {
				missileArt = wFormOne.projectile.Text()
			}
		}
		if wFormTwo.weaponType.Selected() != 0 && wFormTwo.weaponType.Selected() != 1 { // 0 == None, 1 == Normal
			if wFormTwo.projectileSpeed.Text() != "" {
				if missileSpeed != "" {
					missileSpeed += "," + wFormTwo.projectileSpeed.Text()
				} else {
					missileSpeed = wFormTwo.projectileSpeed.Text()
				}
			}
			if wFormTwo.projectileHoming.Checked() {
				if missileHoming != "" {
					missileHoming += ",1"
				} else {
					missileHoming = "1"
				}
			} else {
				if missileHoming != "" {
					missileHoming += ",0"
				} else {
					missileHoming = "0"
				}
			}
			if wFormTwo.projectile.Text() != "" {
				if missileArt != "" {
					missileArt += "," + wFormTwo.projectile.Text()
				} else {
					missileArt = wFormTwo.projectile.Text()
				}
			}
		}
		if missileSpeed != "" {
			newUnitFunc.Missilespeed.SetValid(missileSpeed)
		}
		if missileHoming != "" {
			newUnitFunc.Missilehoming.SetValid(missileHoming)
		}
		if missileArt != "" {
			newUnitFunc.Missileart.SetValid(missileArt)
		}
		if uForm.buttonX.Text() != "" && uForm.buttonY.Text() != "" {
			newUnitFunc.Buttonpos.SetValid(uForm.buttonX.Text() + "," + uForm.buttonY.Text())
		}
		if uForm.icon.Text() != "" {
			newUnitFunc.Art.SetValid(uForm.icon.Text())
		}
		if uForm.tooltip.Text() != "" {
			newUnitFunc.Tip.SetValid(uForm.tooltip.Text())
		}
		if uForm.hotkey.Text() != "" {
			newUnitFunc.Hotkey.SetValid(uForm.hotkey.Text())
		}
		if uForm.description.Text() != "" {
			newUnitFunc.Ubertip.SetValid(uForm.description.Text())
		}
		if dForm.sells.Text() != "" {
			newUnitFunc.Sellitems.SetValid(dForm.sells.Text())
		}
		if dForm.upgradesTo.Text() != "" {
			newUnitFunc.Upgrade.SetValid(dForm.upgradesTo.Text())
		}
		if dForm.trains.Text() != "" {
			newUnitFunc.Trains.SetValid(dForm.trains.Text())
		}

		unitData.UnitID.SetValid(uForm.unitId.Text())
		unitData.Race.SetValid(races[dForm.race.Selected()])
		unitData.Movetp.SetValid(moveTypes[dForm.movementType.Selected()])
		unitData.MoveHeight.SetValid(dForm.flyingHeight.Text())
		unitData.PathTex.SetValid(pathingTextures[uForm.pathingTexture.Selected()])
		unitData.Points.SetValid(dForm.points.Text())

		unitWeapons.UnitWeapID.SetValid(uForm.unitId.Text())
		if wFormOne.enableWeapon.Checked() && wFormTwo.enableWeapon.Checked() {
			unitWeapons.WeapsOn.SetValid("3")
		} else if wFormOne.enableWeapon.Checked() {
			unitWeapons.WeapsOn.SetValid("1")
		} else if wFormTwo.enableWeapon.Checked() {
			unitWeapons.WeapsOn.SetValid("2")
		} else {
			unitWeapons.WeapsOn.SetValid("0")
		}
		unitWeapons.Acquire.SetValid(dForm.acquisition.Text())
		unitWeapons.Targs1.SetValid(targetsToString(wFormOne.targets))
		unitWeapons.SplashTargs1.SetValid(targetsToString(wFormOne.aoeTargets))
		unitWeapons.WeapTp1.SetValid(weaponTypes[wFormOne.weaponType.Selected()])
		unitWeapons.AtkType1.SetValid(attackTypes[wFormOne.attackType.Selected()])
		unitWeapons.BackSw1.SetValid(wFormOne.backswingPoint.Text())
		unitWeapons.Dmgpt1.SetValid(wFormOne.damagePoint.Text())
		unitWeapons.Dmgplus1.SetValid(wFormOne.damageBase.Text())
		unitWeapons.Dice1.SetValid(wFormOne.damageDice.Text())
		unitWeapons.Sides1.SetValid(wFormOne.damageSides.Text())
		unitWeapons.Cool1.SetValid(wFormOne.cooldown.Text())
		unitWeapons.RangeN1.SetValid(wFormOne.weaponRange.Text())
		unitWeapons.Farea1.SetValid(wFormOne.aoeFull.Text())
		unitWeapons.Harea1.SetValid(wFormOne.aoeMedium.Text())
		unitWeapons.Qarea1.SetValid(wFormOne.aoeSmall.Text())
		unitWeapons.Targs2.SetValid(targetsToString(wFormTwo.targets))
		unitWeapons.SplashTargs2.SetValid(targetsToString(wFormTwo.aoeTargets))
		unitWeapons.WeapTp2.SetValid(weaponTypes[wFormTwo.weaponType.Selected()])
		unitWeapons.AtkType2.SetValid(attackTypes[wFormTwo.attackType.Selected()])
		unitWeapons.BackSw2.SetValid(wFormTwo.backswingPoint.Text())
		unitWeapons.Dmgpt2.SetValid(wFormTwo.damagePoint.Text())
		unitWeapons.Dmgplus2.SetValid(wFormTwo.damageBase.Text())
		unitWeapons.Dice2.SetValid(wFormTwo.damageDice.Text())
		unitWeapons.Sides2.SetValid(wFormTwo.damageSides.Text())
		unitWeapons.Cool2.SetValid(wFormTwo.cooldown.Text())
		unitWeapons.RangeN2.SetValid(wFormTwo.weaponRange.Text())
		unitWeapons.Farea2.SetValid(wFormTwo.aoeFull.Text())
		unitWeapons.Harea2.SetValid(wFormTwo.aoeMedium.Text())
		unitWeapons.Qarea2.SetValid(wFormTwo.aoeSmall.Text())

		unitBalance.UnitBalanceID.SetValid(uForm.unitId.Text())
		unitBalance.Level.SetValid(oForm.level.Text())
		// unitBalance.Type.SetValid()
		unitBalance.Goldcost.SetValid(dForm.goldCost.Text())
		unitBalance.Lumbercost.SetValid(dForm.lumberCost.Text())
		unitBalance.Fmade.SetValid(dForm.foodProduction.Text())
		unitBalance.Fused.SetValid(dForm.foodCost.Text())
		unitBalance.HP.SetValid(dForm.health.Text())
		// unitBalance.RegenHP.SetValid(dForm.healthRegen.Text())
		unitBalance.ManaN.SetValid(dForm.mana.Text())
		// unitBalance.RegenMana.SetValid(dForm.manaRegen.Text())
		unitBalance.Def.SetValid(dForm.defense.Text())
		// unitBalance.DefType.SetValid()
		unitBalance.Spd.SetValid(dForm.movementSpeed.Text())
		unitBalance.MinSpd.SetValid(dForm.movementSpeedMinimum.Text())
		unitBalance.MaxSpd.SetValid(dForm.movementSpeedMaximum.Text())
		// unitBalance.Bldtm.SetValid()
		// unitBalance.Reptm.SetValid()
		// unitBalance.Sight.SetValid()
		// unitBalance.Nsight.SetValid()
		// unitBalance.STR.SetValid()
		// unitBalance.AGI.SetValid()
		// unitBalance.INT.SetValid()
		// unitBalance.Upgrades.SetValid()
		if dForm.isBuilding.Checked() {
			unitBalance.Isbldg.SetValid("1")
		} else {
			unitBalance.Isbldg.SetValid("0")
		}

		unitUi.UnitUIID.SetValid(uForm.unitId.Text())
		unitUi.File.SetValid(uForm.model.Text())
		unitUi.UnitSound.SetValid(uForm.soundSet.Text())
		unitUi.Name.SetValid(uForm.name.Text())
		// unitUi.UnitClass.SetValid()
		// unitUi.Special.SetValid()
		if dForm.isCampaign.Checked() {
			unitUi.Campaign.SetValid("1")
		} else {
			unitUi.Campaign.SetValid("1")
		}
		if oForm.dropsItems.Checked() {
			unitUi.DropItems.SetValid("1")
		} else {
			unitUi.DropItems.SetValid("0")
		}
		unitUi.Scale.SetValid(uForm.selectionScale.Text())
		unitUi.ElevPts.SetValid(oForm.elevationSamplePoints.Text())
		unitUi.ElevRad.SetValid(oForm.elevationSampleRadius.Text())
		// unitUi.Walk.SetValid()
		// unitUi.Run.SetValid()
		// unitUi.TeamColor.SetValid()
		// unitUi.CustomTeamColor.SetValid()
		// unitUi.Armor.SetValid()
		unitUi.ModelScale.SetValid(uForm.scalingValue.Text())
		unitUi.Red.SetValid(fmt.Sprint(uForm.red.Value()))
		unitUi.Green.SetValid(fmt.Sprint(uForm.green.Value()))
		unitUi.Blue.SetValid(fmt.Sprint(uForm.blue.Value()))
		// unitUi.UnitShadow.SetValid()
		// unitUi.BuildingShadow.SetValid()
		// unitUi.ShadowW.SetValid()
		// unitUi.ShadowH.SetValid()
		// unitUi.ShadowX.SetValid()
		// unitUi.ShadowY.SetValid()
		// unitUi.ShadowOnWater.SetValid()
		// unitUi.SelCircOnWater.SetValid()

		unitAbilities.UnitAbilID.SetValid(uForm.unitId.Text())
		// unitAbilities.Auto.SetValid()
		unitAbilities.AbilList.SetValid(uForm.abilities.Text())
		// unitAbilities.HeroAbilList.SetValid()

		newSlkUnit.UnitData = unitData
		newSlkUnit.UnitWeapons = unitWeapons
		newSlkUnit.UnitBalance = unitBalance
		newSlkUnit.UnitUI = unitUi
		newSlkUnit.UnitAbilities = unitAbilities
	})
	formVbox.Append(saveUnitButton, false)
	hbox.Append(formVbox, true)

	return hbox
}

func setupUI() {
	mainWindow = ui.NewWindow("WC3 SLK Editor", 640, 480, false)
	mainWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainWindow.Destroy()
		return true
	})

	mainWindow.SetChild(makeBasicControlsPage())
	mainWindow.SetMargined(true)

	mainWindow.Show()
}

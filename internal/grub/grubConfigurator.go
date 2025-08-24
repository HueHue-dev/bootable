package grub

type GrubConfigurator struct {
	builder ICFGBuilder
}

func NewGrubConfigurator(builder ICFGBuilder) *GrubConfigurator {
	return &GrubConfigurator{builder: builder}
}

func (gc *GrubConfigurator) Construct() error {
	if err := gc.builder.createGrubCfgFile(); err != nil {
		return err
	}
	if err := gc.builder.insertHeaderTemplate(); err != nil {
		return err
	}
	if err := gc.builder.insertIsoSpecificTemplates(); err != nil {
		return err
	}
	return gc.builder.GetResult()
}

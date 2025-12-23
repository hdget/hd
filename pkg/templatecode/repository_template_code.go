package templatecode

type RepositoryGenerator interface {
	Gen(destDir, pkgName, repo string) error
}

type repoGeneratorImpl struct{}

const (
	repositoryOperation = `
package {{ .Pkg }}

import (
	"{{ .Pkg }}/autogen/m"
	"{{ .Pkg }}/autogen/pb"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/lib-sqlboiler"
	"github.com/pkg/errors"
)

func (impl {{ .Impl }}) Create(item {{ .BizObject }}) ({{ .ModelObject }}, error) {
	m%s := &m.Company{
		Tid: impl.Tid(),
	}

	err := impl.Copier().CopyForCreate(mCompany, item)
	if err != nil {
		return nil, err
	}

	err = mCompany.Insert(impl.Executor(), boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, "db create company")
	}

	return mCompany, nil
}

func (impl companyTdbImpl) Delete(id int64) error {
	_, err := m.Companies(
		m.CompanyWhere.Tid.EQ(impl.Tid()),
		m.CompanyWhere.ID.EQ(id),
	).DeleteAll(impl.Executor())
	if err != nil {
		return errors.Wrap(err, "db delete company")
	}
	return nil
}

func (impl companyTdbImpl) Edit(item *pb.Company) error {
	if item == nil || item.Id <= 0 {
		return errors.New("invalid param")
	}

	mCompany, err := m.Companies(
		m.CompanyWhere.Tid.EQ(impl.Tid()),
		m.CompanyWhere.ID.EQ(item.Id),
		m.CompanyWhere.Version.EQ(int(item.Version)),
	).One(impl.Executor())
	if err != nil {
		return errors.Wrap(err, "find company")
	}

	// 编辑角色只允许编辑其名称
	err = impl.Copier().Copy(mCompany, item)
	if err != nil {
		return err
	}

	_, err = mCompany.Update(impl.Executor(), boil.Infer())
	if err != nil {
		return errors.Wrap(err, "db edit company")
	}

	return nil
}

func (impl companyTdbImpl) Update(mCompany *m.Company) error {
	mCompany.Version += 1
	_, err := mCompany.Update(impl.Executor(), boil.Infer())
	if err != nil {
		return errors.Wrap(err, "db update company")
	}
	return nil
}

func (impl companyTdbImpl) Get(id int64) (*m.Company, error) {
	return m.Companies(
		m.CompanyWhere.Tid.EQ(impl.Tid()),
		m.CompanyWhere.ID.EQ(id),
	).One(impl.Executor())
}

func (impl companyTdbImpl) Count(filters map[string]string) (int64, error) {
	total, err := m.Companies(
		impl.GetQueryConditions(filters)...,
	).Count(impl.Executor())
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (impl companyTdbImpl) List(filters map[string]string, list ...*protobuf.ListParam) ([]*m.Company, error) {
	mCompanies, err := m.Companies(
		sqlboiler.NewQmBuilder(
			impl.GetQueryConditions(filters)...,
		).Limit(list...).Output()...,
	).All(impl.Executor())
	if err != nil {
		return nil, err
	}

	results := make([]*m.Company, 0)
	for _, mCompany := range mCompanies {
		results = append(results, mCompany)
	}

	return results, nil
}

func (impl companyTdbImpl) GetQueryConditions(filters map[string]string) []qm.QueryMod {
	queryMods := []qm.QueryMod{
		m.CompanyWhere.Tid.EQ(impl.Tid()),
	}

	if v, exists := filters[m.CompanyColumns.Name]; exists {
		queryMods = append(queryMods, m.CompanyWhere.Name.LIKE("%"+v+"%"))
	}

	return queryMods
}

`
)

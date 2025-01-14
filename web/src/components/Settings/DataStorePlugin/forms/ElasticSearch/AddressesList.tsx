import {PlusOutlined} from '@ant-design/icons';
import {Form, Input} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import * as S from '../../DataStorePluginForm.styled';

interface IProps {
  fields: FormListFieldData[];
  add(): void;
  remove(name: number): void;
}

const AddressesList = ({fields, add, remove}: IProps) => {
  return (
    <>
      <S.ItemListContainer>
        {fields.map(({key, name, ...field}, index) => (
          <S.ListItem key={key}>
            <Form.Item key={key} name={[name]} {...{field}} style={{width: '100%'}}>
              <Input placeholder={`Address #${index + 1}`} />
            </Form.Item>
            <S.ItemActionContainer>
              {index !== 0 && (
                <S.DeleteCheckIcon
                  onClick={() => {
                    remove(name);
                  }}
                />
              )}
            </S.ItemActionContainer>
          </S.ListItem>
        ))}
      </S.ItemListContainer>
      <S.AddButton icon={<PlusOutlined />} onClick={() => add()} data-cy="add-assertion-form-add-check">
        Add new
      </S.AddButton>
    </>
  );
};

export default AddressesList;

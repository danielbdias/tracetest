import {useSortable} from '@dnd-kit/sortable';
import {CSS} from '@dnd-kit/utilities';
import {TTest} from 'types/Test.types';
import * as S from './TestsSelectionInput.styled';

interface IProps {
  test: TTest;
  sortableId: string;
  onDelete(sortableId: string): void;
}

const TestItem = ({test, onDelete, sortableId}: IProps) => {
  const {attributes, listeners, setNodeRef, transform, transition} = useSortable({id: sortableId});

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <S.TestItemContainer ref={setNodeRef} style={style} {...attributes}>
      <S.DragHandle {...listeners} />
      <S.TestLink to={`/test/${test.id}`} target="_blank">
        <span>{test.name}</span>
      </S.TestLink>
      <S.DeleteIcon onClick={() => onDelete(sortableId)} />
    </S.TestItemContainer>
  );
};

export default TestItem;

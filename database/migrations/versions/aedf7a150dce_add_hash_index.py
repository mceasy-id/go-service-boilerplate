"""add_hash_index

Revision ID: aedf7a150dce
Revises: 076619c4f740
Create Date: 2023-08-30 08:22:37.379772

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = 'aedf7a150dce'
down_revision = '076619c4f740'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_index('reason_company_id_hash_index', 'product', ['company_id'], unique=False, postgresql_using='hash')
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index('reason_company_id_hash_index', table_name='product', postgresql_using='hash')
    # ### end Alembic commands ###
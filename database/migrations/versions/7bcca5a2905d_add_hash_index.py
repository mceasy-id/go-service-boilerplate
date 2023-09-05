"""add_hash_index

Revision ID: 7bcca5a2905d
Revises: aedf7a150dce
Create Date: 2023-08-30 08:22:51.996266

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '7bcca5a2905d'
down_revision = 'aedf7a150dce'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index('reason_company_id_hash_index', table_name='product')
    op.create_index('product_company_id_hash_index', 'product', ['company_id'], unique=False, postgresql_using='hash')
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index('product_company_id_hash_index', table_name='product', postgresql_using='hash')
    op.create_index('reason_company_id_hash_index', 'product', ['company_id'], unique=False)
    # ### end Alembic commands ###
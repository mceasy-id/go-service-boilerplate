from datetime import datetime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.sql import func, text, expression, schema
from sqlalchemy import (
    DateTime,
    BigInteger,
    Integer,
    Column,
    String,
    Float,
    Boolean,
    ForeignKey
)

Base = declarative_base()

class Product(Base):
    __tablename__ = 'product'

    company_id = Column(BigInteger, nullable=False)
    uuid = Column(UUID(as_uuid=True), primary_key=True)
    name = Column(String, nullable=False)
    description = Column(String, nullable=False)
    price = Column(BigInteger, nullable=True)

    created_on = Column(DateTime(timezone=True), nullable=False, server_default=func.now())
    created_by = Column(String, nullable=False)
    updated_on = Column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_by = Column(String, nullable=False)

    __table_args__  = (
        schema.Index('product_company_id_hash_index',company_id, postgresql_using='hash'),
    )